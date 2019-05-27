package gen

import (
	"errors"
	"fmt"
	"strings"
	"github.com/vibhavp/test-gen/step"
)

var ErrLocatorInvalidType = errors.New("gen: locator missing/invalid type")

var ErrAttributeMissingKey = errors.New("gen: attribute locator missing/invalid key")
var ErrAttributeMissingValue = errors.New("gen: attribute locator missing/invalid value")

var ErrXPathInvalidValue = errors.New("gen: xpath locator missing/invalid value")
var ErrLocatorInvalidPosition = errors.New("gen: invalid locator position")

var ErrStepInvalidLocators = errors.New("gen: missing/invalid locators")
var ErrStepInputInvalidText = errors.New("gen: input step has missing/invalid text")
var ErrStepWaitInvalidValue = errors.New("gen: wait step has missing/invalid value")
var ErrStepWaitInvalidMaxWait = errors.New("gen: wait step has missing/invalid max_wait value")
var ErrStepAssertInvalidValue = errors.New("gen: assertion step has missing/invalid value")
var ErrStepStepUntilInvalid = errors.New("gen: wait_until is missing/invalid")
var ErrStepAssertTypeInvalid = errors.New("gen: assertion type is missing/invalid")


// We convert locators to either XPath or CSS selectors
// (although writing a CSS selector to xpath query would be cool, but unnecessary)

func genXPathFromAttributeLocator(loc *step.Locator) (string, error) {
	if loc.AttributeKey == nil {
		return "", ErrAttributeMissingKey
	}
	if loc.Position != nil {
		return fmt.Sprintf("//*[@%s='%s'][position()=%d]", *loc.AttributeKey, loc.Value, *loc.Position), nil
	}
	return fmt.Sprintf("//*[@%s='%s']", *loc.AttributeKey, loc.Value), nil
}

func (g *GenContext) startTry() {
	g.WriteLine("try:")
	g.currentIndentLevel++
}

func (g *GenContext) endTry() {
	g.currentIndentLevel--
}

func (g *GenContext) genNoElementBlock() {
	g.WriteLine("except NoSuchElementException:")
	g.currentIndentLevel++
	g.WriteLine("elem = None")
	g.currentIndentLevel--
}

func (g *GenContext) genXPathCode(xpath string) {
	xpath = strings.Replace(xpath, `"`, `'`, -1)

	g.WriteLine(fmt.Sprintf(`elem = driver.find_element_by_xpath("%s")`, xpath))
}

func (g *GenContext) genCSSSelectorCode(selector string, position *int) {
	if position == nil {
		g.WriteLine(fmt.Sprintf(`elem = driver.find_element_by_css_selector("%s")`, selector))
	} else {
		g.WriteLine(fmt.Sprintf(`elems = driver.find_elements_by_css_selector("%s")`, selector))
		g.WriteLine(fmt.Sprintf("elem = elems[%d] if len(elems) >= %d else None", *position, *position))
	}
}

func (g *GenContext) genLocatorCode(locators []*step.Locator, throwWhenNotFound bool) error {
	if len(locators) == 0 {
		return nil
	}
	for index, loc := range locators {
		if index > 0 {
			g.WriteLine("if not elem:")
			g.currentIndentLevel++
		}
		g.startTry()
		switch loc.Type {
		case "attribute":
			xpath, err := genXPathFromAttributeLocator(loc)
			if err != nil {
				return err
			}
			g.genXPathCode(xpath)
		case "xpath":
			xpath := loc.Value
			if loc.Position != nil {
				xpath += fmt.Sprintf("[position()=%d]", *loc.Position)
			}
			g.genXPathCode(xpath)
		case "css_selector":
			g.genCSSSelectorCode(loc.Value, loc.Position)
		default:
			return ErrLocatorInvalidType
		}
		g.endTry()
		g.genNoElementBlock()
		if index > 0 {
			g.currentIndentLevel--
		}
	}
	return nil
}

func (g *GenContext) genStep(s *step.Step, funcName string) error {
	g.WriteLine(fmt.Sprintf("def %s(driver):", funcName))
	g.currentIndentLevel++
	if err := g.genLocatorCode(s.Locators, s.Type != "assertion"); err != nil {
		return err
	}
	switch s.Type {
	case "input":
		g.WriteLine(fmt.Sprintf(`elem.send_keys("%s")`, s.InputText))
	case "click":
		if s.Config == nil || s.Config.ClickType == "" {
			g.WriteLine(`elem_type = elem.get_attribute("type")`)
			g.WriteLine(`if elem_type == "submit":`)
			g.currentIndentLevel++
			g.WriteLine(`elem.submit()`)
			g.currentIndentLevel--
			g.WriteLine(`else:`)
			g.currentIndentLevel++
			g.WriteLine(`elem.click()`)
			g.currentIndentLevel--
		} else {
			switch s.Config.ClickType {
			case "submit":
				g.WriteLine(`elem.submit()`)
			case "click":
				g.WriteLine(`elem.click()`)
			}
		}
		
	case "wait":
		switch s.WaitUntil {
		case "url_changed":
			g.WriteLine(fmt.Sprintf(`WebDriverWait(driver, %d).until(EC.url_changes(driver.current_url))`, s.WaitMaxWait))
		case "title_is":
			g.WriteLine(fmt.Sprintf(`WebDriverWait(driver, %d).until(EC.title_is(%s))`, s.WaitMaxWait, s.Value))
		case "title_contains":
			g.WriteLine(fmt.Sprintf(`WebDriverWait(driver, %d).until(EC.title_contains(%s))`, s.WaitMaxWait, s.Value))
		default:
			return ErrStepStepUntilInvalid
		}
	case "assertion":
		switch s.AssertType {
		case "elementNotExists":
			g.WriteLine("assert not elem")
		case "textExists":
			g.WriteLine(fmt.Sprintf(`assert "%s" in elem.text`, s.Value))
		default:
			return ErrStepAssertTypeInvalid
		}
	}

	if s.Config != nil {
		g.WriteLine(fmt.Sprintf(`driver.implicitly_wait(%d)`, s.Config.StepWait))
	}
	g.WriteLine(fmt.Sprintf(`print("%s successful")`, funcName))
	g.WriteLine("\n")
	g.currentIndentLevel--
	return nil
}

func (g *GenContext) GenTest(t *step.Test) error {
	names := make([]string, len(t.Steps))
	
	for index, step := range t.Steps {
		names[index] = fmt.Sprintf("step_%d", index)
		if err := g.genStep(step, names[index]); err != nil {
			return err
		}
	}

	testName := strings.Replace(t.Name, " ", "_", -1)
	testName = strings.ToLower(testName)

	g.WriteLine(fmt.Sprintf(`def test_%s(driver):`, testName))
	g.currentIndentLevel++
	g.WriteLine(fmt.Sprintf(`driver.get("%s")`, t.BaseURL))
	for _, name := range names {
		g.WriteLine(fmt.Sprintf("%s(driver)", name))
	}
	g.currentIndentLevel--

	return nil
}
