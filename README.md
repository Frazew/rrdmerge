# rrdmerge

## Introduction

The RRD (Round-Robin Database) file format remains a fairly simple storage format used by many monitoring tools (most notably LibreNMS). However, `rrdtool` lacks the ability to merge two RRD files together, which can be useful in case of data loss or migrations.

This tool attempts to solve just that: merge two RRD files (or folders) together while trying not to lose any bit of history or override otherwise good data. It also handles rrdcached by flushing (note that you still need access to the filesystem as dumping / restoring over rrdcached is not possible).

## Why not "xxx" other method

There are several scripts or tools on the Internet that aim at merging RRDs. However all the tools found were text-based, using `rrdtool dump`. This leads to poor performance and is error-prone (as illustrated by the countless forum threads seeking help to run them). Because the target use case was to merge thousands of devices, each with hundreds (or thousands) of RRD files, performance was paramount.

Let's not be pretentious, this implementation also has its flaws, but the goal was to have a tool that would both perform reasonably well while also handling errors without **ever** messing up any data.

## Limitations

Because their usage is limited and for lack of real-world RRD files using them, the only supported Consolidation Functions (CF) are the following, any other CF will either fail or lead to (very) unpredictable behavior:

- AVERAGE
- MAXIMUM
- MINIMUM
- LAST

Also note that this tool cannot work over rrdcached solely, it **needs** direct filesystem access because rrdcached does not support dumping and restoring.

## Usage

### Building / Installing

Either use the pre-built binary within this repository `bin/rrdmerge` or build it yourself:

```shell
go mod download
make build
```

Note: the librrd variant (`-tag librrd`) requires librrd to be installed on the system and is dynamically linked. The rrdtool variant requires rrdtool to be installed on the system but has the advantage of being statically linked.

### Running

Using rrdmerge is simple: `rrdmerge -a source_a.rrd -b source_b.rrd`. This will merge the two files and write the result into `source_b.rrd`.

`-a` and `-b` accept both single files and folders. If given a folder, rrdmerge will be run for each RRD file found and copy over to `-b` the ones that are not found (unless `-common` is specified)

The additional flags are:

- `-noSkip` which allows merging files that do not end in ".rrd
- `-t <count>` which specifies the number of copy and merge threads to spawn
- `-d <socket>` which specifies the socket to connect to flush if using rrdcached
- `-s <path>` which allows to specify a base path that should be stripped from B when using rrdcached (absolute paths in rrdcached are only allowed when using a unix domain socket)
- `-dry` makes rrdmerge run without overwriting or copying any files. If `-d` is supplied, calls to rrdcached to flush will still be made

### Behavior

If supplied with folders instead of files, rrdmerge will copy over from A to B the files that do not already exist. If the file exists and has ends in ".rrd" (or the `-noSkip` flag is set), it will attempt to merge them. The merging process is as follows:

- if supplied with a `-d` flag, the "destination" (B) file is flushed with rrdcached
- both RRD files are read into memory, directly parsing the binary file (i.e without using `rrd_dump` or `rrdtool dump`)
- if both files have the same `last update` field, they are not merged
- if the files have a different count of DS, they are not merged
- if A has a more recent `last update` field than B, then the data from A will have precedence over B
- for each RRA, the values are merged by copying the value from A to B (respectively from B to A the previous condition matched) **only if it is not NaN**
- if the difference in time elapsed between A and B for a given RRA is greater than the row count, the RRA is kept intact (such a situation means that the old data we want to merge would already have been overwritten with newer data)
- after merging, the in-memory RRD structures are serialized to XML before getting converted back into binary .rrd using `rrd_restore` or `rrdtool restore`

A few notes:

- rrdmerge assumes that the RRD files have the same RRAs (meaning that RRA #n in A has the same parameters as RRA #n in B)
- additional RRAs are ignored: if A has more RRAs than B, they are not copied over. If B has more RRAs than A, they are left untouched

### Restore behavior

Depending on the build variant (librrd or rrdtool), the behavior for restoring is not the same, which can have unexpected side effects:

- the librrd variant restores the RRD file by writing it to a temporary file in `/tmp/`, like `/tmp/rrdmerge1852119208`
- the rrdtool variant restores the RRD file by piping it into the standard input of `rrdtool restore`

## Use cases

### LibreNMS migration

Consider the following situation: a LibreNMS cluster is using rrdcached with distributed polling. Following a change in rrdcached, the daemon has to be restarted. Because rrdcached is unavailable for a few seconds, the pollers begin caching in their local filesystem the RRD files. When rrdcached comes back up, a few minutes of monitoring data is therefore lost.

Using rrdmerge, we can retrieve these cached RRD files and merge them back into the main storage handled by rrdcached, therefore integrating back the data that would otherwise be lost.

## Testing

Testing is not yet fully implemented as this requires sample (non-production/anonymous) RRD files.

## License

See LICENSE