import { browser, element, by } from 'protractor';

export class OptionsCafeDesktopPage {
  navigateTo() {
    return browser.get('/');
  }

  getParagraphText() {
    return element(by.css('oc-root h1')).getText();
  }
}
