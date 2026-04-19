import { useState, useEffect } from "react";

interface UseCurrentTimeReturn {
  currentTime: Date;
  greeting: string;
  formattedDate: string;
  formattedTime: string;
}

export const useCurrentTime = (updateInterval = 60000): UseCurrentTimeReturn => {
  const [currentTime, setCurrentTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => setCurrentTime(new Date()), updateInterval);
    return () => clearInterval(timer);
  }, [updateInterval]);

  const getGreeting = () => {
    const hour = currentTime.getHours();
    if (hour < 12) return "Доброе утро";
    if (hour < 18) return "Добрый день";
    return "Добрый вечер";
  };

  const formattedDate = currentTime.toLocaleDateString("ru-RU", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  const formattedTime = currentTime.toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
    hourCycle: "h23",
  });

  return {
    currentTime,
    greeting: getGreeting(),
    formattedDate,
    formattedTime,
  };
};
