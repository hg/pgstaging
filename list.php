<?php

require_once 'util.php';

$error = pluck_session('error');
$action = pluck_session('action');

$clusters = json_decode(shell_exec('pg_lsclusters --json'));
usort($clusters, fn($a, $b) => $a->port - $b->port);

function modified_at(string $path): string
{
    $timestamp = filemtime($path);

    if (!$timestamp) {
        return '?';
    }

    $tz = new DateTimeZone('Asia/Almaty');

    $dt = new DateTime;
    $dt->setTimestamp($timestamp);
    $dt->setTimezone($tz);

    return $dt->format('D, d.m.Y, H:i:s');
}
