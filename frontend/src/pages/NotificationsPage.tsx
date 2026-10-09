import { NotificationsPanel } from "../features/NotificationsPanel";

type NotificationsPageProps = {
  notifications: string[];
};

export function NotificationsPage({ notifications }: NotificationsPageProps) {
  return (
    <section className="page-grid single">
      <NotificationsPanel notifications={notifications} />
    </section>
  );
}
