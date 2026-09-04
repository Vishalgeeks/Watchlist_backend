import { ReactNode } from 'react';
import SideNav from './SideNav';
import TopNav from './TopNav';
import BottomNav from './BottomNav';

export default function AppLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex h-screen bg-background">
      <SideNav />
      <div className="flex-1 flex flex-col md:ml-60 min-w-0">
        <TopNav />
        <main className="flex-1 overflow-y-auto pt-16 md:pt-0">
          <div className="p-margin-mobile md:p-margin-desktop max-w-[1600px] mx-auto pb-32 md:pb-margin-desktop">
            {children}
          </div>
        </main>
      </div>
      <BottomNav />
    </div>
  );
}