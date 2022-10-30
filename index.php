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
            <?php if ($action === 'create'): ?>
              База данных создана.
            <?php elseif ($action === 'destroy'): ?>
              База данных удалена.
            <?php elseif ($action === 'stop'): ?>
              База данных остановлена.
            <?php elseif ($action === 'start'): ?>
              База данных запущена.
            <?php endif ?>
        </div>
      <?php endif ?>

    <h1 class="title">Список активных баз</h1>

    <div class="content">
      Не спамим запросами! Скрипты написаны наспех и работают синхронно — браузер может подвисать на минуту. Это
      ожидаемо.
    </div>

    <div class="table-container">
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
                  <form action="modify.php" method="post">
                    <input type="hidden" name="name" value="<?= $cluster->cluster ?>">

                      <?php if ($cluster->running): ?>
                        <input type="hidden" name="action" value="stop">
                        <button type="submit" class="button is-warning is-small" onclick="addLoading()">
                          Остановить
                        </button>
                      <?php else: ?>
                        <input type="hidden" name="action" value="start">
                        <button type="submit" class="button is-primary is-small" onclick="addLoading()">
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
                  <form action="modify.php" method="post">
                    <input type="hidden" name="action" value="destroy">
                    <input type="hidden" name="name" value="<?= $cluster->cluster ?>">

                    <button class="button is-danger is-small" onclick="deleteDb()">
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
    
  </div>
</section>

<section class="section">
  <div class="container">

    <h1 class="title">Создать базу</h1>

    <div class="content">
      <p>
        Допустимые символы: <code>_</code>, <code>abcdefghijklmnopqrstuvwxyz</code>, <code>0123456789</code>.
      </p>

      <p>В качестве имени используйте что-то понятное:</p>

      <ol>
        <li>если это ваша личная песочница — своё имя: <code>divanov</code></li>
        <li>если это выполнение задачи <code>#12345</code> — так и назовите базу: <code>issue_12345</code></li>
        <li>если это работа над обращениями — <code>appeals</code></li>
      </ol>

      <p>
        К остальному подход простой: периодически смотрю список → не понимаю, зачем существует база → удаляю без
        предупреждения.
      </p>
    </div>

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
          <button id="btn-create" class="button is-primary" onclick="addLoading()">
            Создать
          </button>
        </div>
      </div>
    </form>

  </div>
</section>

<section class="section">
  <div class="container">
    <h1 class="title">Дамп</h1>
    <p>Скачать дамп для подъёма на своей копии PostgreSQL можно <a href="dump/sc.dump">здесь</a>.</p>
  </div>
</section>

<script>
  function addLoading() {
    event.target.classList.add("is-loading");
    return true;
  }

  function deleteDb() {
    const ok = confirm("Вы точно хотите удалить базу?");
    if (ok) {
      addLoading();
    }
    return ok;
  }
</script>

</body>
</html>
