import { runCli } from "./cli";

runCli().then((code) => {
  process.exit(code);
});
