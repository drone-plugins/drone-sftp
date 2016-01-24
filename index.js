const Drone = require('drone-node');
const plugin = new Drone.Plugin();

const PromiseSftp = require('promise-sftp');

const path = require('path');
const shelljs = require('shelljs');

const do_upload = function (workspace, vargs) {
  if (vargs.host) {

    var sftp = new PromiseSftp();
    var resolved_files = [].concat.apply([], vargs.files.map((f) => { return shelljs.ls(workspace.path + '/' + f); }));
    vargs.destination_path || (vargs.destination_path = '/');

    //TODO: Support Private Key auth
    sftp.connect({
      host: vargs.host,
      port: vargs.port,
      username: vargs.username,
      password: vargs.password
    }).then(function (greetings) {
      console.log('Connection successful. ' + (greetings || ''));
      
      resolved_files.forEach((file) => {
        var basename = path.basename(file);

        console.log('Uploading ' + file + ' as ' + basename + ' into ' + vargs.destination_path);
        sftp.put(file, path.join(vargs.destination_path, basename))
        .then(function () {

          //TODO: Need to Logout after all the files are uploaded
          console.log('Upload successful')
        }).catch(function (err) {

          console.log('An error happened while trying to push the file ' + file + ': ' + err);
        });
      });
    }).catch(function(err) {

      console.log('Connection to SFTP server failed. ' + err)
    });

  } else {

    console.log("Parameter missing: SFTP server host");
    process.exit(1)
  }
}

plugin.parse().then((params) => {

  // gets build and repository information for
  // the current running build
  const build = params.build;
  const repo  = params.repo;
  const workspace = params.workspace;

  // gets plugin-specific parameters defined in
  // the .drone.yml file
  const vargs = params.vargs;

  vargs.username      || (vargs.username = '');
  vargs.password      || (vargs.password = '');
  vargs.files         || (vargs.files = []);

  do_upload(workspace, vargs);
});
