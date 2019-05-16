package cos

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"path"
	"strings"
)

func stringMap(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func find(vs []string, f func(string) bool) (string, bool) {
	for _, v := range vs {
		ok := f(v)
		if ok {
			return v, ok
		}
	}

	return "", false
}

// try to get _manifests / _layers / _uploads lastIndex
func lastIndexOfImageContent(subPath string) int {
	MANIFESTS := "_manifests"
	LAYERS := "_layers"
	UPLOADS := "_uploads"

	target, ok := find([]string{MANIFESTS, LAYERS, UPLOADS}, func(s string) bool {

		return strings.LastIndex(subPath, s) != -1
	})

	if ok {
		return strings.LastIndex(subPath, target)
	}

	return -1

}

func getPrefix(host string) string {

	chunks := strings.Split(host, ".")

	teamGK := chunks[0]

	return teamGK
}

func getFinalPath(host string, subPath string) string {

	prefix := getPrefix(host)

	repoPrefix := "/docker/registry/v2/repositories/"

	imageContentLastIndex := lastIndexOfImageContent(subPath)

	if strings.HasPrefix(subPath, repoPrefix) && imageContentLastIndex != -1 {
		name := subPath[len(repoPrefix):imageContentLastIndex]

		//log.Printf("[INFO] name:", name)

		nextNameChunks := stringMap(strings.Split(name, "/"), func(c string) string {
			if len(c) > 1 {
				// TODO: 暂时使用一个 hash 算法代替数据库查询 id 的操作
				h := md5.New()
				_, err := io.WriteString(h, c)

				if err != nil {
					log.Printf("[ERROR] get path error:", err)
					return c
				}

				str := fmt.Sprintf("%x", h.Sum(nil))

				return str[0:3]
			}
			return c
		})

		nextSubPath := path.Join(repoPrefix, path.Join(nextNameChunks...), subPath[imageContentLastIndex:])

		nextFullPath := path.Join(prefix, nextSubPath)

		//log.Printf("[INFO] nextFullPath:", nextFullPath)

		return nextFullPath
	}

	return path.Join(prefix, subPath)

}
