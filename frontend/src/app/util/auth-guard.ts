import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { UserStore } from './user.store';

export const authGuard: CanActivateFn = async () => {
  const userStore = inject(UserStore);
  const router = inject(Router);

  const isValid = await userStore.loadAndVerifySession();
  return isValid ? true : router.createUrlTree(['/login']);
};
