# /home/ubuntu/flynn-cleanup.sh
# run as sudo

# install nsenter from S3 (alternatively follow https://github.com/jpetazzo/nsenter)
curl -fsSLo /usr/local/bin/nsenter https://s3.amazonaws.com/lmars.net/nsenter
chmod +x /usr/local/bin/nsenter
 
# leaked mounts are in these directories
mnt_dir="/var/lib/docker/aufs/mnt"
diff_dir="/var/lib/docker/aufs/diff"
 
# list the pids of active LXC processes
lxc_pids=$(ps aux | grep [l]ibvirt_lxc | awk '{print $2}')

# list the mounts which are not for active jobs
leaked=$( find $mnt_dir -maxdepth 1 -name 'tmp-*' -printf '%f\n' | grep -v -f <(flynn-host ps -q | cut -d '-' -f 2-))
 
# for each leaked mount, unmount it from all LXC processes and remove the mnt & diff directories
for name in $leaked; do
  for pid in $lxc_pids; do
     grep -q "${name}" "/proc/${pid}/mounts" &&  nsenter -t "${pid}" -m /bin/umount -fl "${mnt_dir}/${name}"
  done
   rm -rf "${mnt_dir}/${name}" "${diff_dir}/${name}"
done