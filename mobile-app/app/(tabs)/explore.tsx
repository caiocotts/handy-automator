import { StyleSheet } from 'react-native';

import ParallaxScrollView from '@/components/parallax-scroll-view';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { Fonts } from '@/constants/theme';

import React, { useEffect, useState } from 'react';
import { Text, TouchableOpacity, View } from 'react-native';
import { Link } from "expo-router";

export default function TabTwoScreen() {


    const [devices, setDevices] = useState<any[]>([]);
    useEffect(() => {
        getDevices();
    }, []);
    const getDevices = async () => {
        try {
            const response = await fetch("http://localhost:3000/api/device", {
                method: "GET",
                headers: { "Content-Type": "application/json", Authorization: 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzM2MjUyMjIsImlhdCI6MTc3MzYxODAyMiwiaXNzIjoiZGVjaXNpb24tbWFrZXIiLCJzdWIiOiIzem9lZGFCanlmX3gifQ.oXDok282z52spXBQM604aGydOZ18moKqhexOifM8ONg', },
            });
            const data = await response.json();
            setDevices(data.devices ?? []);
        } catch (error) {
            console.error("Error fetching devices:", error);
        }
    };


    return (
        <ParallaxScrollView
            headerBackgroundColor={{ light: '#D0D0D0', dark: '#353636' }}
            headerImage={
                <IconSymbol
                    size={310}
                    color="#808080"
                    name="chevron.left.forwardslash.chevron.right"
                    style={styles.headerImage}
                />
            }>

            <ThemedView style={styles.titleContainer}>
                <ThemedText
                    type="title"
                    style={{
                        fontFamily: Fonts.rounded,
                    }}>
                    Devices
                </ThemedText>
            </ThemedView>
            <ThemedText>Your Home Devices</ThemedText>

            {devices.map((device) => (
                <View key={device.id} style={styles.deviceCard}>
                    <View style={styles.deviceInfo}>
                        <Text style={styles.deviceName}>{device.name ?? 'Unnamed Device'}</Text>
                        {device.type ? <Text style={styles.deviceType}>{device.type}</Text> : null}
                        <Text style={styles.deviceIp}>{device.ip}</Text>
                    </View>
                </View>
            ))}

            {devices.length === 0 && (
                <ThemedText style={{ textAlign: 'center', marginTop: 16 }}>
                    No devices found.
                </ThemedText>
            )}

            <View style={styles.buttonRow}>
                <TouchableOpacity style={styles.button} onPress={getDevices}>
                    <Text style={styles.addButtonText}>Refresh</Text>
                </TouchableOpacity>

                <Link href="/modal" asChild>
                    <TouchableOpacity style={styles.addButton}>
                        <Text style={styles.addButtonText}>Add device</Text>
                    </TouchableOpacity>
                </Link>
            </View>

        </ParallaxScrollView>
    );
}

const styles = StyleSheet.create({
    deviceCard: {
        backgroundColor: '#1E1E2E',
        borderRadius: 12,
        padding: 16,
        marginBottom: 10,
    },
    deviceInfo: {
        gap: 4,
    },
    deviceName: {
        color: '#FFFFFF',
        fontSize: 17,
        fontWeight: '600',
    },
    deviceType: {
        color: '#9BA1A6',
        fontSize: 14,
    },
    deviceIp: {
        color: '#687076',
        fontSize: 13,
    },
    buttonRow: {
        flexDirection: 'row',
        gap: 12,
        marginTop: 8,
    },
    addButton: {
        backgroundColor: '#4e8cc2',
        padding: 10,
        borderRadius: 8,
        alignItems: 'center',
        flex: 1,
    },
    addButtonText: {
        fontSize: 16,
        color: '#FFFFFF',
    },
    button: {
        alignItems: 'center',
        backgroundColor: '#333',
        padding: 10,
        borderRadius: 8,
        flex: 1,
    },
    headerImage: {
        color: '#808080',
        bottom: -90,
        left: -35,
        position: 'absolute',
    },
    titleContainer: {
        flexDirection: 'row',
        gap: 8,
    },
});
