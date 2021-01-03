#!/bin/sh

printf "Command: ${0}"

for arg in "${@}"; do printf " \"${arg}\""; done
if [ "${p1}" != "" ]; then printf " p1=${p1}"; fi
if [ "${p2}" != "" ]; then printf " p2=${p2}"; fi
if [ "${HTTP_METHOD}" != "" ]; then printf " HTTP_METHOD=${HTTP_METHOD}"; fi

exit "${EXITCODE:-0}"
