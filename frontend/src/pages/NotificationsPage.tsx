import { NotificationsPanel } from "../features/NotificationsPanel";
import type { NotificationItem } from "../types";

type NotificationsPageProps = {
  notifications: NotificationItem[];
};

export function NotificationsPage({ notifications }: NotificationsPageProps) {
  return (
    <section className="page-grid single">
      <NotificationsPanel notifications={notifications} />
    </section>
  );
}
