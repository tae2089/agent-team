#!/usr/bin/env node
"use strict";

const fs = require("fs");
const path = require("path");
const childProcess = require("child_process");

const exeName = process.platform === "win32" ? "agent-team-bin.exe" : "agent-team-bin";
const binPath = path.join(__dirname, exeName);

if (!fs.existsSync(binPath)) {
  console.error(
    "agent-team binary is missing. Reinstall @tae2089/agent-team or run `node npm/install.js` from the package directory."
  );
  process.exit(1);
}

const result = childProcess.spawnSync(binPath, process.argv.slice(2), {
  stdio: "inherit",
  windowsHide: false
});

if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}

process.exit(result.status === null ? 1 : result.status);
