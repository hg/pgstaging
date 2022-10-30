<?php

session_start();

function bail(string $error)
{
    $_SESSION['error'] = $error;
    header('Location: /');
    exit(0);
}

function require_post()
{
    if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
        bail('некорректный запрос');
    }
}

function pluck_session(string $key)
{
    $value = $_SESSION[$key];
    unset($_SESSION[$key]);
    return $value;
}

function normalize_name(string $name)
{
    return preg_replace('/[^a-z0-9_]/', '', strtolower($name));
}
