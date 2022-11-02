package crazykfc

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type CrazyKFC struct{ engine.PluginMagic }

var (
	pluginInfo = &CrazyKFC{
		engine.PluginMagic{
			Desc:     "🚀 输入 {kfc骚话} => 获取肯德基疯狂星期四骚话",
			Commands: []string{"kfc骚话"},
		},
	}
	_        = engine.InstallPlugin(pluginInfo)
	sentence []string
)

func (p *CrazyKFC) OnRegister() {
	resp, err := getCrazyKFCSentence()
	if err != nil {
		return
	}
	for i := range resp {
		sentence = append(sentence, resp[i].Text)
	}
}

func (p *CrazyKFC) OnEvent(msg *robot.Message) {
	if msg.MatchTextCommand(pluginInfo.Commands) {
		if len(sentence) > 0 {
			msg.ReplyText(sentence[0])
			sentence = append(sentence[:0], sentence[1:]...)
		} else {
			msg.ReplyText("查询失败，这一定不是bug🤔")
			p.OnRegister()
		}
	}
}

type apiResponse struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

func getCrazyKFCSentence() ([]apiResponse, error) {
	api := "https://fastly.jsdelivr.net/gh/Nthily/KFC-Crazy-Thursday@main/kfc.json"
	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data []apiResponse
	if err := json.Unmarshal(readAll, &data); err != nil {
		return nil, err
	}
	return data, nil
}
