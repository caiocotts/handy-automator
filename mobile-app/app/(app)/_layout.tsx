import { Stack } from 'expo-router';

export default function AppLayout() {
  return (
    <Stack
      screenOptions={{
        headerStyle: { backgroundColor: '#4f46e5' },
        headerTintColor: '#fff',
        headerTitleStyle: { fontWeight: '700' },
      }}
    >
      <Stack.Screen name="workflows" options={{ title: 'Workflows' }} />
      <Stack.Screen name="workflow/[id]" options={{ title: 'Workflow Details' }} />
    </Stack>
  );
}
