#!/bin/sh
set -eu

AGENT_TEAM_STRESS_PROFILE=smoke exec sh "$(dirname -- "$0")/concurrency.sh"
