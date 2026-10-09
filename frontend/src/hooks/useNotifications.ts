import { useEffect, useRef, useState } from "react";
import { WS_BASE } from "../config";

export function useNotifications(token?: string) {
  const [notifications, setNotifications] = useState<string[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!token) {
      wsRef.current?.close();
      wsRef.current = null;
      return;
    }

    const socket = new WebSocket(`${WS_BASE}/ws/notifications?token=${encodeURIComponent(token)}`);
    wsRef.current = socket;

    socket.onopen = () => pushNotification(setNotifications, "Notifications connected");
    socket.onmessage = (event) => pushNotification(setNotifications, event.data);
    socket.onerror = () => pushNotification(setNotifications, "Notification stream error");

    return () => socket.close();
  }, [token]);

  return notifications;
}

function pushNotification(
  setNotifications: React.Dispatch<React.SetStateAction<string[]>>,
  message: string,
) {
  setNotifications((items) => [message, ...items].slice(0, 8));
}
