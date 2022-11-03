# pgstaging

Uses btrfs snapshots under the hood for quickly creating and destroying
throwaway “forks” of your main PostgreSQL instance for integration & end-to-end
testing or experimentation on actual production data.

It was written in a couple of evenings to be used at`$DAYJOB` and thrown over
the fence in the hopes that it might be useful to others as-is (or more likely
my future self). Zero configurability, no security, no promises about anything.

## Installation

Set up a btrfs partition on a real or virtual drive (I'm using the second one
here):

    $ mkdir /opt/dev && cd /opt/dev
    $ fallocate --length 75G disk.img
    $ mkfs.btrfs -m single -d single disk.img

Add a mount to fstab. You might want to leave out compression depending on the
type and size of your data:

    $ cat >>/etc/fstab <<EOF
    /opt/dev/disk.img  /opt/dev/mnt  btrfs  compress=zstd:1,noatime  0  0
    EOF

    $ mount -a

Create a subvolume for the data:

    $ btrfs subvolume create /opt/dev/mnt/base

Set up replication from your main database server. This is a complicated topic,
have a look at the official documentation or blog posts like [this one][repl].

[repl]: https://www.percona.com/blog/2018/11/30/postgresql-streaming-physical-replication-with-slots/

While restoring a backup from master, put your replica's data in that subvolume.
In my case, the subvolume is `/opt/dev/mnt/base`, and the data goes
in `/opt/dev/mnt/base/data`.

Start your replica:

    $ systemctl enable --now postgresql@15-base

Build this project, scp it to the server, install & start it:

    $ CGO_ENABLED=0 go build
    $ scp -C ./db replica:
    $ ssh replica 'sudo ./db install && sudo systemctl start pgstaging'

It will copy itself to `/usr/local/bin` and create a systemd service to manage
the daemon (but won't start right away unless you do it manually).

Open your server's IP or domain in a web browser.
