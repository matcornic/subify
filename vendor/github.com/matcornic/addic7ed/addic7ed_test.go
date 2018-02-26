package addic7ed_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matcornic/addic7ed"
)

func TestAddic7edSearchAllWithGoodShow(t *testing.T) {
	c := addic7ed.New()
	var flagtests = []struct {
		inFile           string
		expectedShowName string
	}{
		{"Shameless.US.S08E11.720p.HDTV.x264-BATV[ettv]", "Shameless (US) - 08x11 - A Gallagher Pedicure"},
		{"This is Us S01E02", "This is Us - 01x02 - The Big Three"},
		{"Transparent 04x09 AMZ.WEB-DL-NTb;WEB.h264-STRIFE", "Transparent - 04x09 - They is on the Way"},
		{"The Big Bang Theory - 06x12 - Web-dl 480p", "The Big Bang Theory - 06x12 - The Egg Salad Equivalency"},
		{"Mr.Robot.S03E09.720p.HDTV.x264-AVS", "Mr. Robot - 03x09 - eps3.8_stage3.torrent"},
		{"Dark.S01E05.720p.WEBRip.x264-STRiFE", "Dark - 01x05 - Truths"},
	}

	for _, test := range flagtests {
		t.Logf("Searching all subtitles for filename %v\n", test.inFile)
		actualShow, err := c.SearchAll(test.inFile)
		assert.NoError(t, err)
		t.Logf("%v: %v", actualShow.Name, actualShow.Subtitles)
		assert.Equal(t, test.expectedShowName, actualShow.Name)
		assert.NotEmpty(t, actualShow.Subtitles)
	}
}

func TestAddic7edAllWithUnknownShow(t *testing.T) {
	c := addic7ed.New()
	fileName := "DoesNotExist"
	t.Logf("Searching show and subtitles with filename: %v\n", fileName)
	show, err := c.SearchAll(fileName)
	if err == nil {
		t.Logf("Show found: %v\n", show.Name)
		t.Logf("Subtitles found: %v\n", show.Subtitles)
		t.Fatal(errors.New("Should not have found show"))
	}

	t.Logf("Got expected error: %v\n", err)
}

func TestAddic7edSearchBestWithGoodShow(t *testing.T) {
	c := addic7ed.New()
	var flagtests = []struct {
		inFile           string
		inLang           string
		expectedShowName string
		expectedVersion  string
	}{
		{"Shameless.US.S08E11.720p.HDTV.x264-BATV[ettv]", "English", "Shameless (US) - 08x11 - A Gallagher Pedicure", "BATV"},
		{"This is Us S01E02", "English", "This is Us - 01x02 - The Big Three", "KILLERS"},
		{"Transparent 04x09 AMZ.WEB-DL-NTb;WEB.h264-STRIFE", "English", "Transparent - 04x09 - They is on the Way", "AMZ.WEB-DL-NTb;WEB.h264-STRIFE"},
		{"The Big Bang Theory - 06x12 - Web-dl 480p", "English", "The Big Bang Theory - 06x12 - The Egg Salad Equivalency", "480p.WEB-DL"},
		{"Mr.Robot.S03E09.720p.HDTV.x264-AVS", "English", "Mr. Robot - 03x09 - eps3.8_stage3.torrent", "AVS-SVA"},
		{"Dark.S01E05.720p.WEBRip.x264-STRiFE", "English", "Dark - 01x05 - Truths", "WEBRip.x264-STRiFE"},
		{"The.Good.Fight.S01E02.1080p.WEB-DL.DD5.1.H264-ViSUM[rartv]", "English", "The Good Fight - 01x02 - First Week", "WEBRip"},
		{"Lethal.Weapon.S01E11.Lawmen.720p.AMZN.WEBRip.DD5.1.x264-NTb[rartv]", "English", "Lethal Weapon - 01x11 - Lawmen", "WEB-DL"},
		{"The.Walking.Dead.S07E11.Hostiles.and.Calamities.720p.AMZN.WEBRip.DD5.1.x264-CasStudio[rartv]", "English", "The Walking Dead - 07x11 - Hostiles and Calamities", "AMZN.WEBRip"},
	}

	for _, test := range flagtests {
		t.Logf("Searching best subtitle for filename %v and language %v\n", test.inFile, test.inLang)
		actualShowName, actualSubtitle, err := c.SearchBest(test.inFile, test.inLang)
		assert.NoError(t, err)
		t.Logf("%v: %v", actualShowName, actualSubtitle)
		assert.Equal(t, test.expectedShowName, actualShowName)
		assert.Equal(t, test.expectedVersion, actualSubtitle.Version)
		assert.Equal(t, test.inLang, actualSubtitle.Language)
	}
}

func TestAddic7edBestWithUnknownShow(t *testing.T) {
	c := addic7ed.New()
	fileName := "DoesNotExist"
	t.Logf("Searching show and subtitles with filename: %v\n", fileName)
	_, _, err := c.SearchBest(fileName, "English")
	assert.Error(t, err)
	t.Logf("Got expected error: %v\n", err)
}

func TestAddic7edBestWithGoodShowButNotSubtitle(t *testing.T) {
	c := addic7ed.New()
	fileName := "Shameless.US.S08E11.720p.HDTV.x264-BATV[ettv]"
	t.Logf("Searching show and subtitles with filename: %v\n", fileName)
	_, _, err := c.SearchBest(fileName, "Japanese")
	assert.Error(t, err)
	t.Logf("Got expected error: %v\n", err)
}

func TestFilterShow(t *testing.T) {
	subs := addic7ed.Subtitles{
		{Version: "A", Language: "French", Link: "http://addic7ed.com/A-good-show"},
		{Version: "A", Language: "French", Link: "http://addic7ed.com/A-good-show-2"},
		{Version: "B", Language: "French", Link: "http://addic7ed.com/A-good-show"},
		{Version: "A", Language: "English", Link: "http://addic7ed.com/A-good-show"},
		{Version: "A", Language: "Italian", Link: "http://addic7ed.com/A-good-show"},
		{Version: "C", Language: "Italian", Link: "http://addic7ed.com/A-good-show"},
	}
	subsWithVersionA := subs.Filter(addic7ed.WithVersion("A"))
	assert.Len(t, subsWithVersionA, 4)
	subsWithVersionAAndLanguageFrench := subsWithVersionA.Filter(addic7ed.WithLanguage("french"))
	assert.Len(t, subsWithVersionAAndLanguageFrench, 2)
}

func TestGroupByVersion(t *testing.T) {
	subs := addic7ed.Subtitles{
		{Version: "A", Language: "French", Link: "http://addic7ed.com/A-good-show"},
		{Version: "A", Language: "French", Link: "http://addic7ed.com/A-good-show-2"},
		{Version: "B", Language: "French", Link: "http://addic7ed.com/A-good-show"},
		{Version: "A", Language: "English", Link: "http://addic7ed.com/A-good-show"},
		{Version: "A", Language: "Italian", Link: "http://addic7ed.com/A-good-show"},
		{Version: "C", Language: "Italian", Link: "http://addic7ed.com/A-good-show"},
	}
	groupByVersion := subs.GroupByVersion()
	assert.Len(t, groupByVersion, 3)
	assert.Len(t, groupByVersion["A"], 4)
	assert.Len(t, groupByVersion["B"], 1)
	assert.Len(t, groupByVersion["C"], 1)
}

func TestGroupByLanguage(t *testing.T) {
	subs := addic7ed.Subtitles{
		{Version: "A", Language: "French", Link: "http://addic7ed.com/A-good-show"},
		{Version: "A", Language: "French", Link: "http://addic7ed.com/A-good-show-2"},
		{Version: "B", Language: "French", Link: "http://addic7ed.com/A-good-show"},
		{Version: "A", Language: "English", Link: "http://addic7ed.com/A-good-show"},
		{Version: "A", Language: "Italian", Link: "http://addic7ed.com/A-good-show"},
		{Version: "C", Language: "Italian", Link: "http://addic7ed.com/A-good-show"},
	}
	groupByLang := subs.GroupByLanguage()
	assert.Len(t, groupByLang, 3)
	assert.Len(t, groupByLang["French"], 3)
	assert.Len(t, groupByLang["English"], 1)
	assert.Len(t, groupByLang["Italian"], 2)
}

func TestFilterRegexp(t *testing.T) {
	subs := addic7ed.Subtitles{
		{Version: "TV-something", Language: "French", Link: "http://addic7ed.com/A-good-show"},
		{Version: "1080p+WAHOO-Team+Dolby-Surround", Language: "French", Link: "http://addic7ed.com/A-good-show-2"},
		{Version: "Another team", Language: "French", Link: "http://addic7ed.com/A-good-show"},
		{Version: "AMZON", Language: "English", Link: "http://addic7ed.com/A-good-show"},
		{Version: "AAMZN.NTb+DEFLATE+ION10", Language: "Italian", Link: "http://addic7ed.com/A-good-show"},
		{Version: "WUT", Language: "Italian", Link: "http://addic7ed.com/A-good-show"},
	}
	regex := regexp.MustCompile(".*WAHOO-Team.*")
	subtitles := subs.Filter(addic7ed.WithVersionRegexp(regex))
	assert.Len(t, subtitles, 1)
}
