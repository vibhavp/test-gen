package step

type Test struct {
	Name string
	Description string
	BaseURL string     `yaml:"base_url"`
	StepWait int       `yaml:"step_wait"`
	Steps []*Step
}

type Step struct {
	Type string
	Description string
	Locators []*Locator
	Config *Config
	Value string     `yaml:"value"`
	
	InputText string  `yaml:"text"`

	WaitUntil string  `yaml:"until"`
	WaitMaxWait int   `yaml:"max_wait"`

	AssertType string `yaml:"assertionType"`
}

const TypeInput = "input"
const TypeClick = "click"
const TypeWait = "wait"

const AssertTypeTextExists = "textExists"
const AssertTypeElementNotExists = "elementNotExists"

type Locator struct {
	Type string
	Value string `yaml:"value"` // for xpath, css_selector and attribute
	Position *int `yaml:"position"`

	AttributeKey *string   `yaml:"key"`	
}

const LocatorTypeXPath = "xpath"
const LoactorTypeAttribute = "attribute"
const LocatorTypeCSSSelector = "css_selector"

type Config struct {
	StepWait int     `yaml:"step_wait"`
	ClickType string `yaml:"click_type"`
}
