import React from 'react';
import type { Route } from '../+types/root';
import Header from '~/routes/layout/Header';
import { Outlet } from 'react-router';

export async function clientLoader() {
  const data = {
    products: ['product 1', 'product 2'],
  };
  return data;
}

export default function Index({ loaderData }: Route.ComponentProps) {
  const products = loaderData;
  console.log(products);
  return (
    <div>
      <Header />
      <main>
        <Outlet />
      </main>
    </div>
  );
}
