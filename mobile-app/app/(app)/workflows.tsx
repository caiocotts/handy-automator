import {
  createWorkflow,
  deleteWorkflow,
  getWorkflows,
  triggerWorkflow,
  Workflow,
} from '../../lib/api';
import { useAuth } from '../../context/AuthContext';
import { useRouter } from 'expo-router';
import { useCallback, useEffect, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  FlatList,
  Modal,
  RefreshControl,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';

export default function WorkflowsScreen() {
  const { logout } = useAuth();
  const router = useRouter();
  const [workflows, setWorkflows] = useState<Workflow[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [triggering, setTriggering] = useState<string | null>(null);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [newName, setNewName] = useState('');
  const [creating, setCreating] = useState(false);

  const fetchWorkflows = useCallback(async () => {
    try {
      const data = await getWorkflows();
      setWorkflows(data);
    } catch (e: any) {
      Alert.alert('Error', e?.message ?? 'Failed to load workflows.');
    }
  }, []);

  useEffect(() => {
    fetchWorkflows().finally(() => setLoading(false));
  }, [fetchWorkflows]);

  async function onRefresh() {
    setRefreshing(true);
    await fetchWorkflows();
    setRefreshing(false);
  }

  async function handleTrigger(id: string, name: string) {
    setTriggering(id);
    try {
      const results = await triggerWorkflow(id);
      const failed = results.filter((r) => !r.ok);
      if (failed.length === 0) {
        Alert.alert('Triggered', `"${name}" triggered successfully.`);
      } else {
        const msgs = failed.map((r) => `• ${r.error ?? r.deviceId}`).join('\n');
        Alert.alert('Partial failure', `Some devices failed:\n${msgs}`);
      }
    } catch (e: any) {
      Alert.alert('Error', e?.message ?? 'Failed to trigger workflow.');
    } finally {
      setTriggering(null);
    }
  }

  async function handleDelete(id: string, name: string) {
    Alert.alert('Delete workflow', `Delete "${name}"?`, [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Delete',
        style: 'destructive',
        onPress: async () => {
          try {
            await deleteWorkflow(id);
            setWorkflows((prev) => prev.filter((w) => w.id !== id));
          } catch (e: any) {
            Alert.alert('Error', e?.message ?? 'Failed to delete workflow.');
          }
        },
      },
    ]);
  }

  async function handleCreate() {
    if (!newName.trim()) {
      Alert.alert('Error', 'Workflow name cannot be empty.');
      return;
    }
    setCreating(true);
    try {
      const w = await createWorkflow(newName.trim());
      setWorkflows((prev) => [...prev, w]);
      setCreateModalVisible(false);
      setNewName('');
    } catch (e: any) {
      Alert.alert('Error', e?.message ?? 'Failed to create workflow.');
    } finally {
      setCreating(false);
    }
  }

  function renderItem({ item }: { item: Workflow }) {
    const isTriggering = triggering === item.id;
    return (
      <TouchableOpacity
        style={styles.card}
        onPress={() => handleTrigger(item.id, item.name)}
        activeOpacity={0.75}
      >
        <View style={styles.cardContent}>
          <View style={styles.cardInfo}>
            <Text style={styles.cardTitle}>{item.name}</Text>
            <Text style={styles.cardSub}>Tap to trigger</Text>
          </View>
          {isTriggering ? (
            <ActivityIndicator color="#4f46e5" />
          ) : (
            <View style={styles.cardActions}>
              <TouchableOpacity
                style={styles.detailsButton}
                onPress={() => router.push(`/(app)/workflow/${item.id}`)}
              >
                <Text style={styles.detailsButtonText}>Details</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={styles.deleteButton}
                onPress={() => handleDelete(item.id, item.name)}
              >
                <Text style={styles.deleteButtonText}>Delete</Text>
              </TouchableOpacity>
            </View>
          )}
        </View>
      </TouchableOpacity>
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
        data={workflows}
        keyExtractor={(item) => item.id}
        renderItem={renderItem}
        contentContainerStyle={styles.list}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <Text style={styles.emptyText}>No workflows yet.</Text>
            <Text style={styles.emptyHint}>Tap + to create one.</Text>
          </View>
        }
      />

      <TouchableOpacity style={styles.fab} onPress={() => setCreateModalVisible(true)}>
        <Text style={styles.fabText}>+</Text>
      </TouchableOpacity>

      <TouchableOpacity style={styles.logoutBtn} onPress={logout}>
        <Text style={styles.logoutText}>Sign out</Text>
      </TouchableOpacity>

      <Modal visible={createModalVisible} animationType="fade" transparent>
        <View style={styles.modalOverlay}>
          <View style={styles.modalCard}>
            <Text style={styles.modalTitle}>New Workflow</Text>
            <TextInput
              style={styles.input}
              placeholder="Workflow name"
              value={newName}
              onChangeText={setNewName}
              autoFocus
            />
            <View style={styles.modalActions}>
              <TouchableOpacity
                style={styles.cancelButton}
                onPress={() => {
                  setCreateModalVisible(false);
                  setNewName('');
                }}
              >
                <Text style={styles.cancelButtonText}>Cancel</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={[styles.createButton, creating && styles.buttonDisabled]}
                onPress={handleCreate}
                disabled={creating}
              >
                {creating ? (
                  <ActivityIndicator color="#fff" />
                ) : (
                  <Text style={styles.createButtonText}>Create</Text>
                )}
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </Modal>
    </View>
  );
}

WorkflowsScreen.options = {
  title: 'Workflows',
};

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f0f4f8' },
  centered: { flex: 1, justifyContent: 'center', alignItems: 'center' },
  list: { padding: 16, paddingBottom: 100 },
  card: {
    backgroundColor: '#fff',
    borderRadius: 14,
    marginBottom: 12,
    shadowColor: '#000',
    shadowOpacity: 0.06,
    shadowRadius: 8,
    shadowOffset: { width: 0, height: 2 },
    elevation: 2,
  },
  cardContent: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 16,
  },
  cardInfo: { flex: 1 },
  cardTitle: { fontSize: 16, fontWeight: '600', color: '#1a1a2e' },
  cardSub: { fontSize: 12, color: '#9ca3af', marginTop: 2 },
  cardActions: { flexDirection: 'row', gap: 8 },
  detailsButton: {
    backgroundColor: '#ede9fe',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 8,
  },
  detailsButtonText: { color: '#4f46e5', fontWeight: '600', fontSize: 13 },
  deleteButton: {
    backgroundColor: '#fee2e2',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 8,
  },
  deleteButtonText: { color: '#dc2626', fontWeight: '600', fontSize: 13 },
  emptyContainer: { alignItems: 'center', marginTop: 60 },
  emptyText: { fontSize: 16, color: '#6b7280' },
  emptyHint: { fontSize: 13, color: '#9ca3af', marginTop: 4 },
  fab: {
    position: 'absolute',
    bottom: 56,
    right: 24,
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: '#4f46e5',
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#4f46e5',
    shadowOpacity: 0.4,
    shadowRadius: 8,
    shadowOffset: { width: 0, height: 4 },
    elevation: 6,
  },
  fabText: { color: '#fff', fontSize: 28, fontWeight: '300', lineHeight: 32 },
  logoutBtn: {
    position: 'absolute',
    bottom: 16,
    right: 24,
  },
  logoutText: { color: '#6b7280', fontSize: 13 },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.4)',
    justifyContent: 'center',
    padding: 24,
  },
  modalCard: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 24,
  },
  modalTitle: { fontSize: 18, fontWeight: '700', color: '#1a1a2e', marginBottom: 16 },
  input: {
    borderWidth: 1,
    borderColor: '#d1d5db',
    borderRadius: 10,
    paddingHorizontal: 14,
    paddingVertical: 12,
    fontSize: 16,
    marginBottom: 16,
    backgroundColor: '#fafafa',
  },
  modalActions: { flexDirection: 'row', gap: 10 },
  cancelButton: {
    flex: 1,
    borderWidth: 1,
    borderColor: '#d1d5db',
    borderRadius: 10,
    paddingVertical: 12,
    alignItems: 'center',
  },
  cancelButtonText: { color: '#6b7280', fontWeight: '600' },
  createButton: {
    flex: 1,
    backgroundColor: '#4f46e5',
    borderRadius: 10,
    paddingVertical: 12,
    alignItems: 'center',
  },
  buttonDisabled: { opacity: 0.6 },
  createButtonText: { color: '#fff', fontWeight: '600' },
});
