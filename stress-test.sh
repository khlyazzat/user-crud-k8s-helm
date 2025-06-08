#!/bin/bash

API_URL="http://arch.homework/otusapp/khlyazzat/v1/user"
COUNT=100
DURATION_MINUTES=5
END_TIME=$((SECONDS + 60 * DURATION_MINUTES))

echo "🚀 Запуск стресс-теста на $DURATION_MINUTES минут..."

while [ $SECONDS -lt $END_TIME ]; do
  echo "🌀 Новый раунд нагрузки..."

  # Создание пользователей
  for i in $(seq 1 $COUNT); do
    curl -s -X POST "$API_URL/create" \
      -H "Content-Type: application/json" \
      -d "{\"name\":\"User$i\",\"email\":\"user${i}x@example.com\",\"age\":25}" > /dev/null &
  done
  wait

  # Обновление
  for i in $(seq 1 $COUNT); do
    curl -s -X PUT "$API_URL/update/$i" \
      -H "Content-Type: application/json" \
      -d "{\"name\":\"UpdatedUser$i\"}" > /dev/null &
  done
  wait

  # Получение
  for i in $(seq 1 $COUNT); do
    curl -s "$API_URL/get/$i" > /dev/null &
  done
  wait

  # Удаление
  for i in $(seq 1 $COUNT); do
    curl -s -X DELETE "$API_URL/delete/$i" > /dev/null &
  done
  wait

  echo "✅ Раунд завершён. Пауза 10 секунд..."
  sleep 10
done

echo "🏁 Стресс-тест завершён!"
