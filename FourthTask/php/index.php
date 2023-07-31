<?php
// Параметры подключения к базе данных
$servername = "localhost";
$username = "root";
$password = "password";
$dbName = "statistic";
$tableName = "t_stat";

// Создание соединения
$conn = new mysqli($servername, $username, $password, $dbName);

// Проверка соединения
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}

// Создание базы данных
$sql = "CREATE DATABASE IF NOT EXISTS $dbName";
if ($conn->query($sql) === TRUE) {
    echo "Database created successfully";
} else {
    echo "Error creating database: " . $conn->error;
}

// Выбор базы данных
mysqli_select_db($conn, $dbName);

// Создание таблицы
$sql = "CREATE TABLE IF NOT EXISTS $tableName (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    path VARCHAR(250) NOT NULL,
    size FLOAT(10, 3) NOT NULL,
    seconds INT NOT NULL,
    dateOfEntry DATE NOT NULL
)";

if ($conn->query($sql) === TRUE) {
    echo "Table created successfully";
} else {
    echo "Error creating table: " . $conn->error;
}

// Обработка POST-запроса
if ($_SERVER["REQUEST_METHOD"] == "POST") {
// Создание соединения
$conn = new mysqli($servername, $username, $password, $dbName);

// Проверка соединения
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}

    // Извлечение данных из POST-запроса
    $path = $_POST["path"];
    $size = (float)$_POST["size"];
    $seconds = (int)$_POST["seconds"];
    $dateOfEntry = $_POST["dateOfEntry"];
    $sql = "INSERT INTO $tableName (path, size, seconds, dateOfEntry)
    VALUES ('$path', '$size', '$seconds', '$dateOfEntry')";

    // выполняем запрос и проверяем на ошибки
    if ($conn->query($sql) === TRUE) {
        echo "New record created successfully";
    } else {
        echo "Error: " . $sql . "<br>" . $conn->error;
    }

// Закрытие соединения с базой данных
$conn->close();
}

// Получение данных из базы данных
$result = mysqli_query($conn, "SELECT * FROM $tableName");

// Создание массива для хранения данных
$data = array();

// Обработка результатов запроса
while ($row = mysqli_fetch_assoc($result)) {
    // Добавление данных в массив
    $data[] = $row;
}
// Вывод содержимого таблицы
$sql = "SELECT id, path, size, seconds, dateOfEntry FROM $tableName";
$result = $conn->query($sql);

if ($result->num_rows > 0) {
    // Вывод данных каждой строки
    while($row = $result->fetch_assoc()) {
        echo "id: " . $row["id"]. " - Name: " . $row["path"]. " " . $row["size"]. " - Email: " . $row["seconds"]. " - Registration Date: " . $row["dateOfEntry"]. "<br>";
    }
} else {
    echo "0 results";
}

// Преобразование массива в JSON и отправка клиенту
header('Content-Type: application/json');
echo json_encode($data);

// Закрытие соединения с базой данных
mysqli_close($conn);
?>