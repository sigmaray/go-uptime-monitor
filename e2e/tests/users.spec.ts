import { test, expect } from '@playwright/test';

async function apiCall(request: any, endpoint: string, data: any) {
  const response = await request.post(`/api/playwright/${endpoint}`, {
    data,
  });
  const text = await response.text();
  expect(response.ok(), `API call failed: ${text}`).toBeTruthy();
  return JSON.parse(text);
}

test.describe('Admin Users Management', () => {
  test.beforeEach(async ({ page, request }) => {
    await apiCall(request, 'clear-table', { table: 'users' });
    await apiCall(request, 'create-user', {
      username: 'admin',
      password: 'password123',
    });

    await page.goto('/login');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button:has-text("Login")');
  });

  test('can navigate to users list and create a new user', async ({ page }) => {
    await page.click('text="Users"');
    await expect(page).toHaveURL('/admin/users');
    
    await page.click('text="Create New User"');
    await expect(page).toHaveURL('/admin/users/new');

    await page.fill('input[name="username"]', 'newguy');
    await page.fill('input[name="password"]', 'secret');
    await page.fill('input[name="confirm_password"]', 'secret');
    await page.click('button:has-text("Create")');

    await expect(page).toHaveURL('/admin/users');
    await expect(page.locator('table')).toContainText('newguy');
  });
});
