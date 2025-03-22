import React from 'react';
import styles from './MainLayout.module.css';

interface MainLayoutProps {
  children: React.ReactNode;
}

export const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  return (
    <div className={styles.mainLayout}>
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.logo}>GuideForge</div>
          <nav className={styles.nav}>
            <ul className={styles.navList}>
              <li className={styles.navItem}><a href="/dashboard">ダッシュボード</a></li>
              <li className={styles.navItem}><a href="/manuals">マニュアル</a></li>
              <li className={styles.navItem}><a href="/profile">プロフィール</a></li>
            </ul>
          </nav>
        </div>
      </header>
      <main className={styles.main}>
        {children}
      </main>
      <footer className={styles.footer}>
        <div className={styles.footerContent}>
          <p>© {new Date().getFullYear()} GuideForge. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}; 