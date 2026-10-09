import { Bell } from "lucide-react";
import type { NotificationItem } from "../types";

type NotificationsPanelProps = {
  notifications: NotificationItem[];
};

export function NotificationsPanel({ notifications }: NotificationsPanelProps) {
  return (
    <section className="panel notifications">
      <div className="panel-title">
        <h2>Notifications</h2>
        <Bell size={18} />
      </div>
      <div className="notification-list">
        {notifications.length === 0 ? (
          <p className="muted">No messages yet.</p>
        ) : (
          notifications.map((item, index) => (
            <article className="notification-item" key={`${item.message}-${item.received_at}-${index}`}>
              <strong>{item.message}</strong>
              <span>{formatReceivedAt(item.received_at)}</span>
            </article>
          ))
        )}
      </div>
    </section>
  );
}

function formatReceivedAt(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  }).format(new Date(value));
}
