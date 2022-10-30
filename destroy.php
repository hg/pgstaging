<?php

require_once 'util.php';

require_post();

$name = normalize_name($_POST['name']);
$prefix = 'dev_';

if (empty($name) ||
    strlen($name) <= strlen($prefix) ||
    strpos($name, $prefix) !== 0
) {
    bail('не указано имя кластера (или оно некорректно)');
}

$destroyed = shell_exec(__DIR__ . "/create destroy $name");

if (!$destroyed) {
    bail('не удалось удалить базу');
}

$_SESSION['action'] = 'deleted';
header('Location: /');
