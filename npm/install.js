#!/usr/bin/env node
"use strict";

const crypto = require("crypto");
const fs = require("fs");
const https = require("https");
const os = require("os");
const path = require("path");

const rootDir = path.resolve(__dirname, "..");
const packageJson = require(path.join(rootDir, "package.json"));
const version = packageJson.version;
const tag = `v${version}`;
const owner = "tae2089";
const repo = "agent-team";
const binDir = path.join(__dirname, "bin");

const platformMap = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows"
};

const archMap = {
  arm64: "arm64",
  x64: "amd64"
};

function fail(message) {
  console.error(message);
  process.exit(1);
}

function target() {
  const goos = platformMap[process.platform];
  const goarch = archMap[process.arch];
  if (!goos || !goarch) {
    fail(`Unsupported platform for agent-team npm install: ${process.platform}/${process.arch}`);
  }

  const ext = goos === "windows" ? ".exe" : "";
  const asset = `agent-team_${tag}_${goos}_${goarch}${ext}`;
  const executable = goos === "windows" ? "agent-team-bin.exe" : "agent-team-bin";

  return {
    asset,
    executable,
    url: `https://github.com/${owner}/${repo}/releases/download/${tag}/${asset}`,
    checksumUrl: `https://github.com/${owner}/${repo}/releases/download/${tag}/SHA256SUMS`
  };
}

function checkOnly() {
  const selected = target();
  if (!version || version.includes("-")) {
    fail("package.json version must be a release version before publishing to npm.");
  }
  if (!packageJson.bin || packageJson.bin["agent-team"] !== "npm/bin/agent-team.js") {
    fail("package.json bin.agent-team must point to npm/bin/agent-team.js.");
  }
  if (!packageJson.publishConfig || packageJson.publishConfig.access !== "public") {
    fail("Scoped npm package must set publishConfig.access to public.");
  }
  console.log(`npm package check passed for ${selected.asset}`);
}

function request(url, redirects = 0) {
  return new Promise((resolve, reject) => {
    https
      .get(
        url,
        {
          headers: {
            "user-agent": `${packageJson.name}/${version}`
          }
        },
        (response) => {
          if (
            response.statusCode >= 300 &&
            response.statusCode < 400 &&
            response.headers.location
          ) {
            response.resume();
            if (redirects >= 5) {
              reject(new Error(`Too many redirects while downloading ${url}`));
              return;
            }
            resolve(request(response.headers.location, redirects + 1));
            return;
          }

          if (response.statusCode !== 200) {
            response.resume();
            reject(new Error(`Download failed (${response.statusCode}) for ${url}`));
            return;
          }

          const chunks = [];
          response.on("data", (chunk) => chunks.push(chunk));
          response.on("end", () => resolve(Buffer.concat(chunks)));
        }
      )
      .on("error", reject);
  });
}

function expectedSha256(sums, asset) {
  for (const line of sums.split(/\r?\n/)) {
    const fields = line.trim().split(/\s+/);
    if (fields.length >= 2 && fields[1] === asset) {
      return fields[0];
    }
  }
  return null;
}

async function install() {
  if (process.env.AGENT_TEAM_SKIP_DOWNLOAD === "1") {
    console.log("Skipping agent-team binary download because AGENT_TEAM_SKIP_DOWNLOAD=1.");
    return;
  }

  const selected = target();
  const destination = path.join(binDir, selected.executable);
  const tempFile = path.join(os.tmpdir(), `${selected.executable}-${process.pid}`);

  fs.mkdirSync(binDir, { recursive: true });

  console.log(`Downloading ${selected.asset}`);
  const checksumData = await request(selected.checksumUrl);
  const expected = expectedSha256(checksumData.toString("utf8"), selected.asset);
  if (!expected) {
    fail(`Could not find ${selected.asset} in SHA256SUMS for ${tag}.`);
  }

  const binary = await request(selected.url);
  const actual = crypto.createHash("sha256").update(binary).digest("hex");
  if (actual !== expected) {
    fail(`Checksum mismatch for ${selected.asset}.`);
  }

  fs.writeFileSync(tempFile, binary, { mode: 0o755 });
  fs.renameSync(tempFile, destination);
  if (process.platform !== "win32") {
    fs.chmodSync(destination, 0o755);
  }
  console.log(`Installed agent-team ${version} for ${process.platform}/${process.arch}`);
}

if (process.argv.includes("--check")) {
  checkOnly();
} else {
  install().catch((error) => fail(error.message));
}
