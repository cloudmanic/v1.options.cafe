import { OptionsCafeDesktopPage } from './app.po';

describe('options-cafe-desktop App', function() {
  let page: OptionsCafeDesktopPage;

  beforeEach(() => {
    page = new OptionsCafeDesktopPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('oc works!');
  });
});
