# Swan Client Tool Guide
[![Made by FilSwan](https://img.shields.io/badge/made%20by-FilSwan-green.svg)](https://www.filswan.com/)
[![Chat on Slack](https://img.shields.io/badge/slack-filswan.slack.com-green.svg)](https://filswan.slack.com)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)

- Join us on our [public Slack channel](https://www.filswan.com/) for news, discussions, and status updates.
- [Check out our medium](https://filswan.medium.com) for the latest posts and announcements.

## Table of Contents
- [Functions](#Functions)
- [Concepts](#Concepts)
- [Prerequisites](#Prerequisites)
- [Installation](#Installation)
- [After Installation](#After-Installation)
- [Configuration](#Configuration)
- [Flowcharts](#Flowcharts)
- [Create Car Files](#Create-Car-Files)
- [Upload Car Files](#Upload-Car-Files)
- [Create A Task](#Create-A-Task)
- [Send deals](#Send-deals)

## Functions
* Generate Car files from downloaded source files with or without Lotus.
* Generate metadata e.g. Car file URI, start epoch, etc. and save them to a metadata CSV file.
* Propose deals based on the metadata CSV file.
* Generate a final CSV file contains deal CIDs and storage provider id for storage provider to import deals.
* Create tasks on Swan Platform.
* Send deal automatically to auto-bid storage providers.

## Concepts

### Task

In swan project, a task can contain multiple offline deals. There are two basic type of tasks:
- Task type
  * **Public Task**: A deal set for open bid. It has 2 types:
    * **Auto-bid** public task: this kind of task will be automatically assigned to a selected storage provider based on reputation system and Market Matcher.
    * **Manual-bid** public task: for this kind of task, after bidder win the bid, the task holder needs to propose the task to the winner. 
  * **Private Task**: It is required to propose deals to a specified storage provider.
- Task status:
  * **Created**: Tasks are created successfully first time on Swan platform for all kinds of tasks.
  * **Assigned**: Tasks have been assigned to storage providers manually by users or automatically by autobid module.
  * **ActionRequired**: Task with autobid mode on,in other words,`bid_mode` set to `1` and `public_deal` set to `true`, have some information missing or invalid:
    - MaxPrice: missing, or is not a valid number
    - FastRetrieval: missing
    - Type: missing, or or not have valid values
    - No offline deals for this task.

    :bell:You need solve the above problems and change the task status to `Created` to participate next run of Market Matcher.
  * **DealSent**: Tasks have been sent to storage providers after tasks being assigned.
  * **ProgressWithFailure**: Some and not all of deals of the task have been sent to storage providers after tasks being assigned.

### Offline Deal

- The size of an offline deal can be up to 64 GB.
- Every step will generate a JSON file which contains file(s) description like below:
```json
[
 {
  "Uuid": "261ac5ae-cfbb-4ae2-a924-3361075b0e60",
  "SourceFileName": "test3.txt",
  "SourceFilePath": "[source_file_dir]/test3.txt",
  "SourceFileMd5": "17d6f25c72392fc0895b835fa3e3cf52",
  "SourceFileSize": 43857488,
  "CarFileName": "test3.txt.car",
  "CarFilePath": "[car_file_dir]/test3.txt.car",
  "CarFileMd5": "9eb7d54ac1ed8d3927b21a4dcd64a9eb",
  "CarFileUrl": "http://[IP]:[PORT]/ipfs/Qmb7TMcABYnnM47dznCPxpJKPf9LmD1Yh2EdZGvXi2824C",
  "CarFileSize": 12402995,
  "DealCid": "bafyreiccgalsj2a3wtrxygcxpp2hfq3h2fwafh63wcld3uq5hakyimpura",
  "DataCid": "bafykbzacecpuzwmiaxc2u4r5bb7p3ukkhotmkfw4mfv3un6huvk6ctugowikq",
  "PieceCid": "baga6ea4seaqjcip2xh265h2pucvwxv7seeawm4gfksfua4zsbb24zujplzsukja",
  "MinerFid": "[miner_fid]",
  "StartEpoch": 1266686,
  "SourceId": 2
 }
]
```
- This JSON file generated in each step is used in its next step and can be used to rebuild the graph in the future.
- Uuid is generated for future index purpose.

## Prerequisites

- Lotus node

## Installation
### Option:one:  **Prebuilt package**: See [release assets](https://github.com/filswan/go-swan-client/releases)
```shell
wget https://github.com/filswan/go-swan-client/releases/download/release-0.1.0/install.sh
chmod +x ./install.sh
./install.sh
```

### Option:two:  Source Code
:bell:**go 1.16+** is required
```shell
git clone https://github.com/filswan/go-swan-client.git
cd go-swan-client
git checkout <release_branch>
chmod +x ./build_from_source.sh
./build_from_source.sh
```

## After Installation
- The binary file `go-swan-client` is generated under `./build` directory, you need to switch to it.
```shell
cd build
```
- Before executing, you should check your configuration in `~/.swan/client/config.toml` to ensure it is right.
```shell
vi ~/.swan/client/config.toml
```

## Configuration

### [lotus]
- **client_api_url**:  Url of lotus client web api, such as: **http://[ip]:[port]/rpc/v0**, generally the [port] is **1234**
- **client_access_token**:  Access token of lotus client web api. It should have admin access right. You can get it from your lotus node machine using command `lotus auth create-token --perm admin`. See [Obtaining Tokens](https://docs.filecoin.io/build/lotus/api-tokens/#obtaining-tokens)

### [main]

- **api_url**: Swan API address. For Swan production, it is "https://api.filswan.com". It can be ignored if offline_mode is set to true in `[sender]` section
- :bangbang:**api_key**: Your api key. Acquire from [Swan Platform](https://www.filswan.com/) -> "My Profile"->"Developer Settings". It can be ignored if offline_mode is set to true in `[sender]` section.
- :bangbang:**access_token**: Your access token. Acquire from [Swan Platform](https://www.filswan.com/) -> "My Profile"->"Developer Settings". It can be ignored if offline_mode is set to true in `[sender]` section.
- :bangbang:**storage_server_type**: "ipfs server" or "web server"

### [web-server]

- **download_url_prefix**: web server url prefix, such as: `https://[ip]:[port]/download`. Store car files for downloading by storage provider, car file url will be `[download_url_prefix]/[filename]`
### [ipfs-server]

- **download_url_prefix**: ipfs server url prefix, such as: `http://[ip]:[port]/ipfs`. Store car files for downloading by storage provider, car file url will be `[download_url_prefix]/[filename]`
- **upload_url**: ipfs server url for uploading file, such as `http://[ip]:[port]/api/v0/add?stream-channels=true&pin=true`

### [sender]

- **bid_mode**: [0/1] Default 1, which is auto-bid mod and it means swan will automatically allocate storage provider for it, while 0 is manual-bid mode and it needs to be bidded manually by storage providers.
- **offline_mode**: [true/false] Default false. If it is set to true, you will not be able to create Swan task on filswan.com, but you can still Car Files, CSV and JSON files for sending deals
- **output_dir**: When you do not set -out-dir option in your command, it is used as the default output directory for saving generated car files, CSV and JSON files. Should be absolute path and you need have access right to this folder or to create it.
- **public_deal**: [true/false] Whether deals in the tasks are public deals
- **verified_deal**: [true/false] Whether deals in this task are going to be sent as verified
- **fast_retrieval**: [true/false] Indicates that data should be available for fast retrieval
- **generate_md5**: [true/false] Whether to generate md5 for each car file, note: this is a resource consuming action.
- **skip_confirmation**: [true/false] Whether to skip manual confirmation of each deal before sending
- **wallet**:  Wallet used for sending offline deals
- **max_price**: Max price willing to pay per GiB/epoch for offline deal
- **start_epoch_hours**: start_epoch for deals in hours from current time
- **expired_days**: expected completion days for storage provider sealing data
- **gocar_file_size_limit**: go car file size limit in bytes
- **duration**: expressed in blocks (1 block is equivalent to 30s). Default value is 1512000, that is 525 days.

## Flowcharts

### Option:one:
- **Conditions:** `[sender].public_deal=true` and `[sender].bid_mode=0`, see [Configuration](#Configuration)
<img src="http://yuml.me/diagram/plain/activity/(start)->(Create Car Files)->(Upload Car Files)->(Create Public Manual-Bid Task)->(Send Deals)->(end)" >


- Task status change process in this option:

[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggVERcbiAgICBBW0NyZWF0ZWRdIC0tPiBCW0Fzc2lnbmVkXSIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In0sInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0)](https://mermaid-js.github.io/mermaid-live-editor/edit/#eyJjb2RlIjoiZ3JhcGggVERcbiAgICBBW0NyZWF0ZWRdIC0tPiBCW0Fzc2lnbmVkXSIsIm1lcm1haWQiOiJ7XG4gIFwidGhlbWVcIjogXCJkZWZhdWx0XCJcbn0iLCJ1cGRhdGVFZGl0b3IiOmZhbHNlLCJhdXRvU3luYyI6dHJ1ZSwidXBkYXRlRGlhZ3JhbSI6ZmFsc2V9)


### Option:two:
- **Conditions:** `[sender].public_deal=true` and `[sender].bid_mode=1`, see [Configuration](#Configuration)
<img src="http://yuml.me/diagram/plain/activity/(start)->(Create Car Files)->(Upload Car Files)->(Create Public Auto-Bid Task)->(Send Auto-Bid Deals)->(end)" >


- Task status change process in this option. Below are some possibilities:

[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggVERcbiAgICBBW0NyZWF0ZWRdIC0tPiBDe01lZXQgRjF9XG4gICAgQyAtLT58Tm98IERbQWN0aW9uUmVxdWlyZWRdXG4gICAgQyAtLT58WWVzfCBFe01lZXQgRjJ9XG4gICAgRSAtLT4gfE5PfCBBXG4gICAgRSAtLT4gfFlFU3wgRltBc3NpZ25lZF0gLS0-R1tTZW5kIERlYWxdIC0tPiBIe0RlYWwgU2VudD99XG4gICAgSCAtLT4gfDB8IEZcbiAgICBIIC0tPiB8QUxMfCBJW0RlYWxTZW50XVxuICAgIEggLS0-IHxTT01FfCBKW1Byb2dyZXNzV2l0aEZhaWx1cmVdIiwibWVybWFpZCI6eyJ0aGVtZSI6ImRlZmF1bHQifSwidXBkYXRlRWRpdG9yIjpmYWxzZSwiYXV0b1N5bmMiOnRydWUsInVwZGF0ZURpYWdyYW0iOmZhbHNlfQ)](https://mermaid-js.github.io/mermaid-live-editor/edit/#eyJjb2RlIjoiZ3JhcGggVERcbiAgICBBW0NyZWF0ZWRdIC0tPiBDe01lZXQgRjF9XG4gICAgQyAtLT58Tm98IERbQWN0aW9uUmVxdWlyZWRdXG4gICAgQyAtLT58WWVzfCBFe01lZXQgRjJ9XG4gICAgRSAtLT4gfE5PfCBBXG4gICAgRSAtLT4gfFlFU3wgRltBc3NpZ25lZF0gLS0-R1tTZW5kIERlYWxdIC0tPiBIe0RlYWwgU2VudD99XG4gICAgSCAtLT4gfDB8IEZcbiAgICBIIC0tPiB8QUxMfCBJW0RlYWxTZW50XVxuICAgIEggLS0-IHxTT01FfCBKW1Byb2dyZXNzV2l0aEZhaWx1cmVdIiwibWVybWFpZCI6IntcbiAgXCJ0aGVtZVwiOiBcImRlZmF1bHRcIlxufSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0)

- F1: all conditions for auto-bid that a task and its offline deals must match when run market matcher
- F2: market matcher can select a miner that can meet all conditions of auto-bid for this task and its offline deals

### Option:three:
- **Conditions:** `[sender].public_deal=false` and `[sender].bid_mode=0`, see [Configuration](#Configuration)
<img src="http://yuml.me/diagram/plain/activity/(start)->(Create Car Files)->(Upload Car Files)->(Create Private Task)->(end)" >

- Task status change process in this option:

[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggVERcbiAgICBBW0NyZWF0ZWRdIiwibWVybWFpZCI6eyJ0aGVtZSI6ImRlZmF1bHQifSwidXBkYXRlRWRpdG9yIjpmYWxzZSwiYXV0b1N5bmMiOnRydWUsInVwZGF0ZURpYWdyYW0iOmZhbHNlfQ)](https://mermaid-js.github.io/mermaid-live-editor/edit/#eyJjb2RlIjoiZ3JhcGggVERcbiAgICBBW0NyZWF0ZWRdIiwibWVybWFpZCI6IntcbiAgXCJ0aGVtZVwiOiBcImRlZmF1bHRcIlxufSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0)

## Create Car Files
:bell: The input dir and out dir should only be absolute one.

:bell: This step is necessary for both public and private tasks. You can choose one of the following 2 options.

### Option:one: By lotus web json rpc api
```shell
./go-swan-client car -input-dir [input_files_dir] -out-dir [car_files_output_dir]
```
**Command parameters used in this step:**
- -input-dir(Required): The directory where the source files reside in.
- -out-dir(optional): Car files and metadata files will be generated into this directory. When omitted, use `[sender].output_dir` in [Configuration](#Configuration)

**Configurations used in this step:**
- [lotus].client_api_url, see [Configuration](#Configuration)
- [lotus].client_access_token, see [Configuration](#Configuration)
- [sender].output_dir, only used when -out-dir is omitted in command, see [Configuration](#Configuration)
- [sender].generate_md5, when it is true, then generate md5 for source files and car files, see [Configuration](#Configuration)

**Files generated after this step:**
- car.csv: contains information for both source files and car files
- car.json: contains information for both source files and car files
- [source-file-name].car: each source file has a related car file

### Option:two: By graphsplit api
```shell
./go-swan-client gocar -input-dir [input_files_dir] -out-dir [car_files_output_dir]
```
**Command parameters used in this step:**
- -input-dir(Required): The directory where the source files reside in.
- -out-dir(optional): Car files and metadata files will be generated into this directory. When omitted, use `[sender].output_dir` in [Configuration](#Configuration)

**Configurations used in this step:**
- [lotus].client_api_url, see [Configuration](#Configuration)
- [lotus].client_access_token, see [Configuration](#Configuration)
- [sender].gocar_file_size_limit, see [Configuration](#Configuration)
- [sender].output_dir, only used when -out-dir is omitted in command, see [Configuration](#Configuration)
- [sender].generate_md5, when it is true, then generate md5 for source files and car files, see [Configuration](#Configuration)

**Files generated after this step:**
- car.csv: contains information for both source files and car files
- car.json: contains information for both source files and car files
- [source-file-name].car: each source file has one or more related car file(s) according to its size and `[sender].gocar_file_size_limit`

Credits should be given to filedrive-team. More information can be found in https://github.com/filedrive-team/go-graphsplit.

## Upload Car Files
:bell: It is required to upload car files to file server, either to web server or to ipfs server.

### Option:one: To a web-server manually
```shell
no go-swan-client subcommand should be executed
```
**Configurations used in this step:**
- [main].storage_server_type, it should be set to `web server`, see [Configuration](#Configuration)

### Option:two: To a local ipfs server
```shell
./go-swan-client upload -input-dir [input_file_dir]
```
**Command parameters used in this step:**
- -input-dir(Required): The directory where the car files and metadata files reside in. Metadata files will be used and updated after car files uploaded.

**Configurations used in this step:**
- [main].storage_server_type, it should be set to `ipfs server` see [Configuration](#Configuration)
- [ipfs_server].download_url_prefix, see [Configuration](#Configuration)
- [sender].output_dir, only used when -out-dir is omitted in command, see [Configuration](#Configuration)

**Files updated after this step:**
- car.csv: car file url will be updated on the original one
- car.json: car file url will be updated on the original one

## Create A Task
:bell: This step is necessary for both public and private tasks. You can choose one of the following 3 options.

### Option:one: Private Task
- **Conditions:** `[sender].public_deal=false`, see [Configuration](#Configuration)
```shell
./go-swan-client task -input-dir [car_files_dir] -out-dir [output_files_dir] -miner [Storage_provider_id] -dataset [curated_dataset] -description [description]
```
**Command parameters used in this step:**
- -input-dir(Required): Input directory where the generated car files and metadata files reside in.
- -out-dir(optional): Metadata files and swan task file will be generated to this directory. When ommitted, use default `[send].output_dir`, see [Configuration](#Configuration)
- -miner(Required): Storage provider Id you want to send deal to
- -dataset(optional): The curated dataset from which the car files are generated
- -description(optional): Details to better describe the data and confine the task or anything the storage provider needs to be informed.

**Configurations used in this step:**
- [sender].public_deal, see [Configuration](#Configuration)
- [sender].bid_mode, see [Configuration](#Configuration)
- [sender].verified_deal, see [Configuration](#Configuration)
- [sender].offline_mode, see [Configuration](#Configuration)
- [sender].fast_retrieval, see [Configuration](#Configuration)
- [sender].max_price, see [Configuration](#Configuration)
- [sender].start_epoch_hours, see [Configuration](#Configuration)
- [sender].expire_days, see [Configuration](#Configuration)
- [sender].generate_md5, when it is true and there is no md5 in car.json, then generate md5 for source files and car files, see [Configuration](#Configuration)
- [sender].wallet, see [Configuration](#Configuration)
- [sender].skip_confirmation, see [Configuration](#Configuration)
- [sender].duration, see [Configuration](#Configuration)
- [main].storage_server_type, see [Configuration](#Configuration)
- [main].api_url, see [Configuration](#Configuration)
- [main].api_key, see [Configuration](#Configuration)
- [main].access_token, see [Configuration](#Configuration)
- [lotus].api_url, see [Configuration](#Configuration)
- [lotus].access_token, see [Configuration](#Configuration)
- [sender].output_dir, only used when -out-dir is omitted in command, see [Configuration](#Configuration)

**Files generated after this step:**
- [task-name].csv: a CSV generated for posting a task and its offline deals on Swan platform or transferring to storage providers directly for offline import
- [task-name]-metadata.csv: contains more contents used for review, uuid will be added based upon car.csv generated in last step
- [task-name]-metadata.json: contains more content for creating proposal in the next step, uuid will be added based upon car.csv generated in last step

### Option:two: Public and Auto-Bid Task
- **Conditions:** `[sender].public_deal=true` and `[sender].bid_mode=1`, see [Configuration](#Configuration)
```shell
./go-swan-client task -input-dir [car_files_dir] -out-dir [output_files_dir] -dataset [curated_dataset] -description [description]
```
**Command parameters used in this step:**
- -input-dir(Required): Input directory where the generated car files and metadata files reside in.
- -out-dir(optional): Metadata files and swan task file will be generated to this directory. When ommitted, use default `[send].output_dir`, see [Configuration](#Configuration)
- -dataset(optional): The curated dataset from which the car files are generated
- -description(optional): Details to better describe the data and confine the task or anything the storage provider needs to be informed.

**Configurations used in this step:**
- [sender].public_deal, see [Configuration](#Configuration)
- [sender].bid_mode, see [Configuration](#Configuration)
- [sender].verified_deal, see [Configuration](#Configuration)
- [sender].offline_mode, see [Configuration](#Configuration)
- [sender].fast_retrieval, see [Configuration](#Configuration)
- [sender].max_price, see [Configuration](#Configuration)
- [sender].start_epoch_hours, see [Configuration](#Configuration)
- [sender].expire_days, see [Configuration](#Configuration)
- [sender].generate_md5, when it is true and there is no md5 in car.json, then generate md5 for source files and car files, see [Configuration](#Configuration)
- [sender].duration, see [Configuration](#Configuration)
- [main].storage_server_type, see [Configuration](#Configuration)
- [main].api_url, see [Configuration](#Configuration)
- [main].api_key, see [Configuration](#Configuration)
- [main].access_token, see [Configuration](#Configuration)
- [sender].output_dir, only used when -out-dir is omitted in command, see [Configuration](#Configuration)

**Files generated after this step:**
- [task-name].csv: a CSV generated for posting a task and its offline deals on Swan platform or transferring to storage providers directly for offline import
- [task-name]-metadata.csv: contains more contents used for review, uuid will be added based upon car.csv generated in last step
- [task-name]-metadata.json: contains more content for creating proposal in the next step, uuid will be added based upon car.csv generated in last step

### Option:three: Public and Manual-Bid Task
- **Conditions:** `[sender].public_deal=true` and `[sender].bid_mode=0`, see [Configuration](#Configuration)
```shell
./go-swan-client task -input-dir [car_files_dir] -out-dir [output_files_dir] -dataset [curated_dataset] -description [description]
```
**Command parameters used in this step:**
- -input-dir(Required): Input directory where the generated car files and metadata files reside in.
- -out-dir(optional): Metadata files and swan task file will be generated to this directory. When ommitted, use default `[send].output_dir`, see [Configuration](#Configuration)
- -dataset(optional): The curated dataset from which the car files are generated
- -description(optional): Details to better describe the data and confine the task or anything the storage provider needs to be informed.

**Configurations used in this step:**
- [sender].public_deal, see [Configuration](#Configuration)
- [sender].bid_mode, see [Configuration](#Configuration)
- [sender].verified_deal, see [Configuration](#Configuration)
- [sender].offline_mode, see [Configuration](#Configuration)
- [sender].fast_retrieval, see [Configuration](#Configuration)
- [sender].max_price, see [Configuration](#Configuration)
- [sender].start_epoch_hours, see [Configuration](#Configuration)
- [sender].expire_days, see [Configuration](#Configuration)
- [sender].generate_md5, when it is true and there is no md5 in car.json, then generate md5 for source files and car files, see [Configuration](#Configuration)
- [sender].duration, see [Configuration](#Configuration)
- [main].storage_server_type, see [Configuration](#Configuration)
- [main].api_url, see [Configuration](#Configuration)
- [main].api_key, see [Configuration](#Configuration)
- [main].access_token, see [Configuration](#Configuration)
- [sender].output_dir, only used when -out-dir is omitted in command, see [Configuration](#Configuration)

**Files generated after this step:**
- [task-name].csv: a CSV generated for posting a task and its offline deals on Swan platform or transferring to storage providers directly for offline import
- [task-name]-metadata.csv: contains more contents used for review, uuid will be added based upon car.csv generated in last step
- [task-name]-metadata.json: contains more content for creating proposal in the next step, uuid will be added based upon car.csv generated in last step

## Send deals
:bell: The input dir and out dir should only be absolute one.

:bell: This step is only necessary for public tasks. You can choose one of the following 2 options according to your task bid_mode.

### Option:one: Manual deal

**Conditions:**
- `task can be found by uuid from swan platform`
- `task.public_deal=true`
- `task.bid_mode=0`

**Note** This step is used only for public deals, since for private deals, the step `Create A Task` includes sending deals.
```shell
./go-swan-client deal -json [task-name-metadata.json] -out-dir [output_files_dir] -miner [storage_provider_id]
```
**Command parameters used in this step:**
- -json(Required): File path to the metadata json file. Mandatory metadata fields: source_file_size, car_file_url, data_cid, piece_cid
- -out-dir(optional): Swan deal final metadata files will be generated to the given directory. When ommitted, use default: `[sender].output_dir`, see [Configuration](#Configuration)
- -miner(Required): Target storage provider id, e.g f01276

**Configurations used in this step:**
- [sender].wallet, see [Configuration](#Configuration)
- [sender].verified_deal, see [Configuration](#Configuration)
- [sender].fast_retrieval, see [Configuration](#Configuration)
- [sender].start_epoch_hours, see [Configuration](#Configuration)
- [sender].skip_confirmation, see [Configuration](#Configuration)
- [sender].max_price, see [Configuration](#Configuration)
- [sender].duration, see [Configuration](#Configuration)
- [main].api_url, see [Configuration](#Configuration)
- [main].api_key, see [Configuration](#Configuration)
- [main].access_token, see [Configuration](#Configuration)
- [sender].output_dir, only used when -out-dir is omitted in command, see [Configuration](#Configuration)

**Files generated after this step:**
- [task-name].csv: a CSV generated for updating offline deal status and fill deal CID for offline deals
- [task-name]-deals.csv: deal CID added based on [task-name]-metadata.csv generated on next step
- [task-name]-deals.json: deal CID added based on [task-name]-metadata.json generated on next step

### Option:two: Auto-bid deal

- **Conditions:** `[sender].public_deal=true` and `[sender].bid_mode=1`, see [Configuration](#Configuration)
- **Note** After swan allocated a miner to a task, the client needs to sending auto-bid deal using the information submitted to swan in step `Create A Task`
```shell
./go-swan-client auto -out-dir [output_files_dir]
```
**Command parameters used in this step:**
- -miner(Required): Target storage provider id, e.g f01276

**Configurations used in this step:**
- [sender].wallet, see [Configuration](#Configuration)
- [sender].max_price, see [Configuration](#Configuration)
- [main].api_url, see [Configuration](#Configuration)
- [main].api_key, see [Configuration](#Configuration)
- [main].access_token, see [Configuration](#Configuration)
- [sender].output_dir, only used when -out-dir is omitted in command, see [Configuration](#Configuration)

**Files generated for each task after this step:**
- [task-name].csv: a CSV generated for updating task status and fill deal CID for offline deals
- [task-name]-deals.csv: deal CID added based on [task-name]-metadata.csv generated on next step
- [task-name]-deals.json: deal CID added based on [task-name]-metadata.json generated on next step
