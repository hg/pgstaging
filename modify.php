<?php

require_once 'util.php';

require_post();

$name = normalize_name($_POST['name']);
$prefix = 'dev_';

if (empty($name) ||
    strlen($name) <= strlen($prefix) ||
    strpos($name, $prefix) !== 0
) {
    bail('не указано имя (или оно некорректно)');
}

$action = $_POST['action'];

switch ($action) {
    case 'destroy':
    case 'stop':
    case 'start':
        break;

    default:
        bail("неизвестное действие $action");
}

exec(__DIR__ . "/create $action $name", $output, $result);

if ($result !== 0) {
    bail('не удалось выполнить действие');
}

$_SESSION['action'] = $action;
header('Location: /');
