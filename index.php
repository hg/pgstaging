<!doctype html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport"
        content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Тестовые базы данных</title>
  <link rel="stylesheet" href="bulma.min.css">
</head>
<body>

<?php require_once 'list.php' ?>

<section class="section">
  <div class="container">

      <?php if ($error): ?>
        <div class="notification is-danger">
          Не удалось выполнить действие: <?= $error ?>
        </div>
      <?php endif ?>

      <?php if ($action): ?>
        <div class="notification is-primary">
            <?php if ($action === 'created'): ?>
              База данных создана.
            <?php elseif ($action === 'deleted'): ?>
              База данных удалена.
            <?php endif ?>
        </div>
      <?php endif ?>

    <h1 class="title">Список активных баз</h1>

    <table class="table is-bordered is-striped">
      <thead>
      <tr>
        <th>Имя</th>
        <th>Порт</th>
        <th>Юзер</th>
        <th>Пароль</th>
        <th>Запущен</th>
        <th>Изменено</th>
        <th>Данные</th>
        <th>Действия</th>
      </tr>
      </thead>

      <tbody>
      <?php foreach ($clusters as $cluster): ?>
          <?php $devDb = strpos($cluster->cluster, 'dev_') === 0; ?>
        <tr>
          <td><?= $cluster->cluster ?></td>
          <td><?= $cluster->port ?></td>
          <td>sc</td>
          <td>sc</td>
          <td>
              <?php if ($devDb): ?>
                <form action="lifecycle.php" method="post">
                  <input type="hidden" name="name" value="<?= $cluster->name ?>">

                    <?php if ($cluster->running): ?>
                      <input type="hidden" name="action" value="stop">
                      <button type="submit" class="button is-warning is-small">
                        Остановить
                      </button>
                    <?php else: ?>
                      <input type="hidden" name="action" value="start">
                      <button type="submit" class="button is-primary is-small">
                        Запустить
                      </button>
                    <?php endif ?>
                </form>
              <?php else: ?>
                  <?= $cluster->running ? '✅' : '❌' ?>
              <?php endif ?>
          </td>
          <td>
              <?= modified_at($cluster->pgdata) ?>
          </td>
          <td><?= $cluster->pgdata ?></td>
          <td>
              <?php if ($devDb): ?>
                <form action="destroy.php" method="post">
                  <input type="hidden" name="name" value="<?= $cluster->name ?>">

                  <button class="button is-danger is-small" onclick="return confirm('Вы точно хотите удалить базу?')">
                    Удалить
                  </button>
                </form>
              <?php endif ?>
          </td>
        </tr>
      <?php endforeach ?>
      </tbody>
    </table>

  </div>
</section>

<section class="section">
  <div class="container">

    <h1 class="title">Создать базу</h1>

    <form action="create.php" method="post">
      <div class="field">
        <label for="name" class="label">Название</label>
        <div class="control">
          <input id="name" type="text" class="input" minlength="1" maxlength="32" name="name"
                 placeholder="Короткое название латиницей/цифрами">
        </div>
      </div>

      <div class="field">
        <div class="control">
          <button id="btn-create" class="button is-primary">
            Создать
          </button>
        </div>
      </div>
    </form>

  </div>
</section>

</body>
</html>
