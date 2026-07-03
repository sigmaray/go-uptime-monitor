import { test, expect } from '@playwright/test';

// Helper to make API requests to our Playwright API
async function apiCall(request: any, endpoint: string, data: any) {
  const response = await request.post(`/api/playwright/${endpoint}`, {
    data,
  });
  const text = await response.text();
  expect(response.ok(), `API call failed: ${text}`).toBeTruthy();
  return JSON.parse(text);
}

test.describe('Authentication', () => {
  test.beforeEach(async ({ request }) => {
    // Clear users table before each test
    await apiCall(request, 'clear-table', { table: 'users' });
  });

  test('should login successfully with valid credentials', async ({ page, request }) => {
    // Create a test user via API
    await apiCall(request, 'create-user', {
      username: 'testuser',
      password: 'password123',
    });

    // Go to login page
    await page.goto('/login');
    
    // Fill credentials
    await page.fill('input[name="username"]', 'testuser');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button:has-text("Login")');

    // Should redirect to admin dashboard
    await expect(page).toHaveURL('/admin/');
    await expect(page.locator('h2')).toContainText('Admin Dashboard');
    
    // Execute SQL directly as a test for the sql endpoint
    const res = await apiCall(request, 'sql', {
      query: 'SELECT COUNT(*) as count FROM users',
    });
    expect(res.status).toBe('ok');
  });

  test('should fail login with invalid credentials', async ({ page }) => {
    await page.goto('/login');
    
    await page.fill('input[name="username"]', 'wronguser');
    await page.fill('input[name="password"]', 'wrongpass');
    await page.click('button:has-text("Login")');

    // Should stay on login page and show error
    await expect(page).toHaveURL('/login');
    // Assuming there's some error message displayed, but let's just check we're still on login page
    await expect(page.locator('button[type="submit"]')).toBeVisible();
  });
});
