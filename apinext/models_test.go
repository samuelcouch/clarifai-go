package main

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var goodModelzBody string = `
<h1>Model Config</h1>Specify the left hand key (without :something) as the model arg in requests which maps to a backend spire_conf name.
<pre>{'celeb-v1.1:facedet': ['51169_sorta2'],
 'celeb-v1.1:facedetrec': ['51169_sorta2'],
 'default:embed': ['4691_sim_no_tree_center1'],
 'default:facedet': ['51169_sorta2'],
 'default:facedetrec': ['51169_sorta2'],
 'default:tag': ['30065_sorta2'],
 'general-v1.1:embed': ['4691_sim_no_tree_center1'],
 'general-v1.1:facedet': ['51169_sorta2'],
 'general-v1.1:facedetrec': ['51169_sorta2'],
 'general-v1.1:tag': ['30065_sorta2']
 }</pre><h1>Backend Map</h1>Map from each spire_conf name to a list host:port where spires are detected.<pre>{'20348_sorta2': [10.0.2.212:1234],
 '24023_center1': [10.0.4.108:1234],
 '25293_center1': [10.0.0.188:1232],
 '30065_sorta2': [10.0.4.108:1233, 10.0.0.188:1234, 10.0.5.112:1231, 10.0.2.212:1233, 10.0.2.212:1232],
 '40727_sorta2': [10.0.4.109:1232],
 '41443_sorta2': [10.0.4.108:1232],
 '4691_sim_no_tree_center1': [10.0.4.108:1230, 10.0.0.188:1233, 10.0.5.112:1232, 10.0.2.212:1230],
 '51169_sorta2': [10.0.4.108:1231, 10.0.0.188:1231, 10.0.5.112:1230, 10.0.2.212:1231],
 '51173_sorta2': [10.0.4.109:1234, 10.0.0.188:1230, 10.0.5.112:1234],
 '70736_sorta2': [10.0.4.109:1231],
 '80893_sorta2': [10.0.4.109:1233, 10.0.4.109:1230, 10.0.5.112:1233]}</pre>
 `

type ByName []ModelInfo

func (a ByName) Len() int { return len(a) }

func (a ByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}
func (a ByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Order-independent test that elements are equal.
func EqualsIgnoreOrder(a, b []ModelInfo) bool {
	if len(a) != len(b) {
		return false
	}
	var namesa []string
	for _, v := range a {
		namesa = append(namesa, v.Name)
	}
	var namesb []string
	for _, v := range b {
		namesb = append(namesb, v.Name)
	}
	if !HasSameElements(namesa, namesb) {
		return false
	}
	m := make(map[string][]string)
	for _, v := range a {
		m[v.Name] = v.SupportedOps
	}
	for _, v := range b {
		if _, ok := m[v.Name]; !ok {
			return false
		}
		if !HasSameElements(v.SupportedOps, m[v.Name]) {
			return false
		}
	}
	return true
}

func HasSameElements(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	return reflect.DeepEqual(a, b)
}

func TestGoodModelz(t *testing.T) {
	// Inject response using a function object.  See http://openmymind.net/Dependency-Injection-In-Go/
	getModelzResponse = func() (*goquery.Document, error) {
		return goquery.NewDocumentFromReader(strings.NewReader(goodModelzBody))
	}

	resp, _ := getModelsFromModelz()
	expectedModels := []ModelInfo{
		ModelInfo{"celeb-v1.1", []string{"facedet", "facedetrec"}},
		ModelInfo{"default", []string{"embed", "facedet", "facedetrec", "tag"}},
		ModelInfo{"general-v1.1", []string{"embed", "facedet", "facedetrec", "tag"}},
	}
	sort.Sort(ByName(resp.Models))
	sort.Sort(ByName(expectedModels))
	if !EqualsIgnoreOrder(resp.Models, expectedModels) {
		t.Errorf("Unexpected results: %v vs %v", resp.Models, expectedModels)
	}
}
