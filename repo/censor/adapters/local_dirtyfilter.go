package adapters

import (
	"fmt"
	filter "github.com/antlinker/go-dirtyfilter"
	"github.com/antlinker/go-dirtyfilter/store"
	"go_another_chatgpt/repo/censor"
	"io/ioutil"
	"path"
	"strings"
)

type LocalDirtyFilter struct {
	filterManage *filter.DirtyManager
	invalid      []rune
}

const InvalidWords = " ~!@#$%^&*()_-+=?<>.—,，。/\\\\|《》？;:：'‘；“\""

func NewLocalDirtyFilter() (r *LocalDirtyFilter) {
	r = new(LocalDirtyFilter)

	//_, filename, _, _ := runtime.Caller(0)
	// The ".." may change depending on you folder structure
	//dir := path.Join(path.Dir(filename), "../../..")
	dir := "."
	wordDir := path.Join(dir, "resources/sensitive_words")

	fs, err := ioutil.ReadDir(wordDir)
	if err != nil {
		panic(err)
	}

	words := []string{}
	for _, file := range fs {
		if !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
		bys, err := ioutil.ReadFile(path.Join(wordDir, file.Name()))
		if err != nil {
			panic(fmt.Errorf("reading %v failed, err=%w", file.Name(), err))
		}
		fileWords := strings.Split(strings.TrimSpace(string(bys)), "\n")
		fileWords2 := []string{}
		for _, word := range fileWords {
			word = strings.TrimSpace(word)
			if word == "" {
				continue
			}
			fileWords2 = append(fileWords2, word)
		}
		words = append(words, fileWords2...)
	}
	//words = []string{"习近平"}
	//log.Printf("%#v\n", words)

	memStore, err := store.NewMemoryStore(store.MemoryConfig{
		DataSource: words,
	})
	if err != nil {
		panic(err)
	}
	r.filterManage = filter.NewDirtyManager(memStore)

	r.invalid = []rune{}
	for _, w := range InvalidWords {
		//fmt.Printf("%c\n", w)
		r.invalid = append(r.invalid, w)
	}
	return
}

func (r *LocalDirtyFilter) MakeTextAuditing(id, text string) (result *censor.TextAuditingResult, err error) {
	//fmt.Printf("text=%#v, invalid=%#v\n", text, invalid)
	res, err := r.filterManage.Filter().Replace(text, '*')
	if err != nil {
		panic(err)
	}
	//fmt.Printf("结果 ：%#v\n", res)
	result = &censor.TextAuditingResult{}

	result.FilteredText = res
	result.Safe = text == res
	return
}
