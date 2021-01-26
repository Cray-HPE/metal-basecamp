@Library("dst-shared@release/shasta-1.4") _
dockerBuildPipeline {
 githubPushRepo = "Cray-HPE/basecamp"
 githubPushBranches = "(release/.*|main)"
 app = "basecamp"
 name = "basecamp"
 description = "Datasource for serving cloud-init metadata."
 dockerfile = "Dockerfile"
 repository = "metal"
 imagePrefix = "cloud"
 product = "csm"
 slackNotification = ["", "", false, true, true, true]
 lintScript = "noop.sh"
 unitTestScript = "noop.sh"
}
