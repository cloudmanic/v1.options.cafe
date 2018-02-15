// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.

export const environment = {
  production: false,
  app_server: "http://127.0.0.1:7080",
  ws_server: "ws://127.0.0.1:7080",
  default_page_size: 25     
};
