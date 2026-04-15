import {
  associateWorkflowDevices,
  Device,
  getDevices,
  getWorkflow,
  Workflow,
} from '../../../lib/api';
import { useLocalSearchParams, useNavigation } from 'expo-router';
import { useCallback, useEffect, useLayoutEffect, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  FlatList,
  RefreshControl,
  StyleSheet,
  Switch,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';

export default function WorkflowDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const navigation = useNavigation();

  const [workflow, setWorkflow] = useState<Workflow | null>(null);
  const [allDevices, setAllDevices] = useState<Device[]>([]);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [saving, setSaving] = useState(false);

  const loadData = useCallback(async () => {
    try {
      const [wf, devices] = await Promise.all([getWorkflow(id), getDevices()]);
      setWorkflow(wf);
      setAllDevices(devices);
      const associated = new Set((wf.devices ?? []).map((d) => d.id));
      setSelectedIds(associated);
    } catch (e: any) {
      Alert.alert('Error', e?.message ?? 'Failed to load workflow details.');
    }
  }, [id]);

  useEffect(() => {
    loadData().finally(() => setLoading(false));
  }, [loadData]);

  useLayoutEffect(() => {
    if (workflow) {
      navigation.setOptions({ title: workflow.name });
    }
  }, [workflow, navigation]);

  async function onRefresh() {
    setRefreshing(true);
    await loadData();
    setRefreshing(false);
  }

  function toggleDevice(deviceId: string) {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(deviceId)) {
        next.delete(deviceId);
      } else {
        next.add(deviceId);
      }
      return next;
    });
  }

  async function handleSave() {
    setSaving(true);
    try {
      await associateWorkflowDevices(id, Array.from(selectedIds));
      Alert.alert('Saved', 'Device associations updated.');
    } catch (e: any) {
      Alert.alert('Error', e?.message ?? 'Failed to save device associations.');
    } finally {
      setSaving(false);
    }
  }

  function renderDevice({ item }: { item: Device }) {
    const isSelected = selectedIds.has(item.id);
    return (
      <View style={styles.deviceRow}>
        <View style={styles.deviceInfo}>
          <Text style={styles.deviceName}>{item.name ?? item.hostname}</Text>
          <Text style={styles.deviceHostname}>{item.hostname}</Text>
        </View>
        <Switch
          value={isSelected}
          onValueChange={() => toggleDevice(item.id)}
          trackColor={{ false: '#d1d5db', true: '#a5b4fc' }}
          thumbColor={isSelected ? '#4f46e5' : '#9ca3af'}
        />
      </View>
    );
  }

  if (loading) {
    return (
      <View style={styles.centered}>
        <ActivityIndicator size="large" color="#4f46e5" />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <FlatList
        data={allDevices}
        keyExtractor={(item) => item.id}
        renderItem={renderDevice}
        contentContainerStyle={styles.list}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        ListHeaderComponent={
          <Text style={styles.sectionHeader}>Associated Devices</Text>
        }
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <Text style={styles.emptyText}>No devices registered.</Text>
          </View>
        }
      />

      <View style={styles.footer}>
        <TouchableOpacity
          style={[styles.saveButton, saving && styles.buttonDisabled]}
          onPress={handleSave}
          disabled={saving}
        >
          {saving ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text style={styles.saveButtonText}>Save Device Associations</Text>
          )}
        </TouchableOpacity>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f0f4f8' },
  centered: { flex: 1, justifyContent: 'center', alignItems: 'center' },
  list: { padding: 16, paddingBottom: 100 },
  sectionHeader: {
    fontSize: 13,
    fontWeight: '600',
    color: '#6b7280',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
    marginBottom: 10,
  },
  deviceRow: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 14,
    marginBottom: 10,
    shadowColor: '#000',
    shadowOpacity: 0.05,
    shadowRadius: 6,
    shadowOffset: { width: 0, height: 2 },
    elevation: 1,
  },
  deviceInfo: { flex: 1 },
  deviceName: { fontSize: 15, fontWeight: '600', color: '#1a1a2e' },
  deviceHostname: { fontSize: 12, color: '#9ca3af', marginTop: 2 },
  emptyContainer: { alignItems: 'center', marginTop: 40 },
  emptyText: { fontSize: 15, color: '#6b7280' },
  footer: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    padding: 16,
    backgroundColor: '#f0f4f8',
  },
  saveButton: {
    backgroundColor: '#4f46e5',
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: 'center',
  },
  buttonDisabled: { opacity: 0.6 },
  saveButtonText: { color: '#fff', fontSize: 16, fontWeight: '600' },
});
