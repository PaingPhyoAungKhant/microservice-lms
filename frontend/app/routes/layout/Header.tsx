import React from 'react';

import HamburgerMenu from '~/presentation/components/layout/HamburgerMenu';
export default function Header() {
  return (
    <header className="w-full px-4 py-4 md:px-30 md:py-8">
      <div className="flex flex-row items-center justify-between">
        <div className="h-17 w-66">
          <img className="h-full w-full" src="/logo/logo-no-background.png" alt="ASTO Logo" />
        </div>
        <HamburgerMenu />
      </div>
    </header>
  );
}
