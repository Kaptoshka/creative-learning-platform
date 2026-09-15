import React, { createContext, useState, useEffect, useContext, useCallback } from 'react';

const OfflineQueueContext = createContext();

export const OfflineQueueProvider = ({ children }) => {
  const [queue, setQueue] = useState(() => {
    // 1. При загрузке восстанавливаем очередь из localStorage
    const saved = localStorage.getItem('offlineQueue');
    return saved ? JSON.parse(saved) : [];
  });

  const [isOnline, setIsOnline] = useState(navigator.onLine);

  // 2. Слушаем состояние сети
  useEffect(() => {
    const handleOnline = () => {
      setIsOnline(true);
      processQueue(); // Как только появился интернет — пробуем отправить
    };
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, [queue]);

  // 3. Сохраняем очередь при каждом изменении
  useEffect(() => {
    localStorage.setItem('offlineQueue', JSON.stringify(queue));
  }, [queue]);

  // 4. Функция добавления в очередь
  const addToQueue = (item) => {
    const newItem = {
      id: crypto.randomUUID(),
      timestamp: Date.now(),
      attempts: 0,
      ...item, // { url, method, body, timeout }
    };

    setQueue((prev) => [...prev, newItem]);

    // Если интернет есть, пробуем отправить сразу
    if (navigator.onLine) {
      setTimeout(() => processQueue(), 0);
    }
  };

  // 5. Логика обработки очереди
  const processQueue = useCallback(async () => {
    if (!navigator.onLine || queue.length === 0) return;

    // Копируем очередь, чтобы не мутировать стейт напрямую во время итерации
    const currentQueue = [...queue];
    const failedIds = [];
    const successfulIds = [];

    for (const item of currentQueue) {
      // Проверка таймаута (если задан)
      if (item.timeout && Date.now() - item.timestamp > item.timeout) {
        console.warn('Request timed out:', item);
        successfulIds.push(item.id);
        continue;
      }

      try {
        const response = await fetch(item.url, {
          method: item.method || 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(item.body),
        });

        if (!response.ok && response.status < 500) {
           throw new Error('Client Error');
        }

        successfulIds.push(item.id);
      } catch (error) {
        console.error('Failed to send item:', item.id, error);
        // Если сеть упала в процессе — прерываем цикл
        if (!navigator.onLine) break;

        // Если ошибка 500 или другая — оставляем в очереди (failedIds можно использовать для счетчика попыток)
        failedIds.push(item.id);
      }
    }

    // Удаляем успешные (и просроченные) из стейта
    if (successfulIds.length > 0) {
      setQueue((prev) => prev.filter((item) => !successfulIds.includes(item.id)));
    }
  }, [queue]);

  return (
    <OfflineQueueContext.Provider value={{ addToQueue, queue, isOnline }}>
      {children}
    </OfflineQueueContext.Provider>
  );
};

export const useOfflineQueue = () => useContext(OfflineQueueContext);
