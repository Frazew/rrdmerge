#!/bin/bash

DRY_RUN=0
DRY_OPT=""
STRIP_PATH=""
LIMIT=0

ARGUMENTS=()
while [ $# -gt 0 ]; do
	[ "$1" == "--dry-run" ] && { DRY_RUN=1; DRY_OPT="-dry"; echo "Warning: dry run"; shift; continue; }
	[ "$1" == "--limit" ] && { shift; LIMIT="$1"; echo "Warning: limiting to $LIMIT devices"; shift; continue; }
	[ "$1" == "--strip" ] && { shift; STRIP_PATH="$1"; echo "Warning: stripping destination path with $1"; shift; continue; }
	ARGUMENTS+=("$1")
	shift
done

set -- "${ARGUMENTS[@]}"

[ $# -lt 3 ] && { echo "Usage: $0 [--dry-run] [--limit n] [--strip path] <folder_a> <folder_b> <rrdcached_sock>"; exit 1; }
[ -z "$NMS_TOKEN" ] && { echo "Please set the NMS_TOKEN environment variable: export NMS_TOKEN=\"yourtoken\""; exit 1; }
[ -z "$NMS_URL" ] && { echo "Please set the NMS_URL environment variable: export NMS_URL=\"nms.example.com\""; exit 1; }
[ -d "$1" ] || { echo "$1 is not a valid directory"; exit 1; }
[ -d "$2" ] || { echo "$2 is not a valid directory"; exit 1; }

i=0
total=$(find "$1" -maxdepth 1 -mindepth 1 -type d | wc -l)

while IFS= read -r -d $'\0' file; do
	device=$(basename "$file")
	[ ! "$LIMIT" -eq 0 ] && [ $i -eq "$LIMIT" ] && { echo "Stopping at $i devices according to limit $LIMIT"; break; }
	[ $((i % 50)) -eq 0 ] && echo "Processed $i/$total devices"

	[ -d "$2/$device" ] || { echo "Skipping device $device as it no longer exists"; i=$((i+1)); continue; }
	echo "### $device ###" >> /tmp/rrd_merge.log
	
	is_disabled=0
	if curl --silent -H "X-Auth-Token: $NMS_TOKEN" "https://$NMS_URL/api/v0/devices/$device" | grep -q 'disabled": 1'; then
		is_disabled=1
	fi

	if ! ( 
		([ $DRY_RUN -eq 0 ] && [ $is_disabled -eq 0 ] && curl -X PATCH -d '{"field": "disabled", "data": 1}' -H "X-Auth-Token: $NMS_TOKEN" "https://$NMS_URL/api/v0/devices/$device");
		./rrdmerge -a "$1/$device" -b "$2/$device" -t 4 -d "$3" -s "$STRIP_OPT" "$DRY_OPT";
		([ $DRY_RUN -eq 0 ] && [ $is_disabled -eq 0 ] && curl -X PATCH -d '{"field": "disabled", "data": 0}' -H "X-Auth-Token: $NMS_TOKEN" "https://$NMS_URL/api/v0/devices/$device")
	) >> /tmp/rrd_merge.log 2>&1; then echo "Device $device failed"; fi

	echo >> /tmp/rrd_merge.log
	i=$((i+1))
done < <(find "$1" -maxdepth 1 -mindepth 1 -type d -print0)