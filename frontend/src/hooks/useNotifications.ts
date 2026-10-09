import { useEffect, useRef, useState } from "react";
import { WS_BASE } from "../config";
import type { NotificationItem } from "../types";

export function useNotifications(token?: string) {
  const [notifications, setNotifications] = useState<NotificationItem[]>([]);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectsRef = useRef(0);

  useEffect(() => {
    if (!token) {
      wsRef.current?.close();
      wsRef.current = null;
      reconnectsRef.current = 0;
      return;
    }

    const encodedToken = encodeURIComponent(token);
    let stopped = false;
    let reconnectTimer: number | undefined;

    function connect() {
      const socket = new WebSocket(`${WS_BASE}/ws/notifications?token=${encodedToken}`);
      wsRef.current = socket;

      socket.onopen = () => {
        reconnectsRef.current = 0;
        pushNotification(setNotifications, "Notifications connected");
      };
      socket.onmessage = (event) => pushNotification(setNotifications, event.data);
      socket.onerror = () => pushNotification(setNotifications, "Notification stream error");
      socket.onclose = () => {
        if (stopped) return;
        const attempts = reconnectsRef.current + 1;
        reconnectsRef.current = attempts;
        if (attempts > 5) {
          pushNotification(setNotifications, "Notification stream disconnected");
          return;
        }

        const delay = Math.min(1000 * attempts, 5000);
        pushNotification(setNotifications, `Reconnecting notifications in ${delay / 1000}s`);
        reconnectTimer = window.setTimeout(connect, delay);
      };
    }

    connect();

    return () => {
      stopped = true;
      if (reconnectTimer) window.clearTimeout(reconnectTimer);
      wsRef.current?.close();
    };
  }, [token]);

  return notifications;
}

function pushNotification(
  setNotifications: React.Dispatch<React.SetStateAction<NotificationItem[]>>,
  message: string,
) {
  setNotifications((items) => [{ message, received_at: new Date().toISOString() }, ...items].slice(0, 8));
}
