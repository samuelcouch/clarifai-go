package main

import (
	//"fmt"
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
 'general-v1.1:tag': ['30065_sorta2'],
 'general-v1.2:embed': ['51173_sorta2'],
 'general-v1.2:embed_tag': ['51173_sorta2'],
 'general-v1.2:tag': ['51173_sorta2'],
 'general-v1.3:embed': ['80893_sorta2'],
 'general-v1.3:embed_tag': ['80893_sorta2'],
 'general-v1.3:tag': ['80893_sorta2']
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

func TestGoodModelz(t *testing.T) {
	// Inject response using a function object.  See http://openmymind.net/Dependency-Injection-In-Go/
	getModelzResponse = func() (*goquery.Document, error) {
		return goquery.NewDocumentFromReader(strings.NewReader(goodModelzBody))
	}

	resp, _ := getModelsFromModelz()
	expectedModels := []string{"celeb-v1.1", "default", "general-v1.1", "general-v1.2", "general-v1.3"}
	if len(resp.Models) != len(expectedModels) {
		t.Errorf("Wrong number of models parsed: %d vs %d", len(resp.Models), len(expectedModels))
	}
}
