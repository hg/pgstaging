<?php

require_once 'util.php';

require_post();

$name = normalize_name($_POST['name']);

if (empty($name) || strlen($name) > 32) {
    bail("некорректное имя '$name'");
}

$created = shell_exec(__DIR__ . "/create create dev_$name");

if (!$created) {
    bail('не удалось клонировать базу');
}

$_SESSION['action'] = 'create';
header('Location: /');
