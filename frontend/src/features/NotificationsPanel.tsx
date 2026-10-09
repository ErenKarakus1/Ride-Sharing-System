import { Bell } from "lucide-react";

type NotificationsPanelProps = {
  notifications: string[];
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
          notifications.map((item, index) => <p key={`${item}-${index}`}>{item}</p>)
        )}
      </div>
    </section>
  );
}
