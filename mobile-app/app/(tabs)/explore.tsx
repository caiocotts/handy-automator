import { Image } from 'expo-image';
import { Alert, Platform, StyleSheet } from 'react-native';

import { Collapsible } from '@/components/ui/collapsible';
import { ExternalLink } from '@/components/external-link';
import ParallaxScrollView from '@/components/parallax-scroll-view';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { Fonts } from '@/constants/theme';
import { Button } from '@react-navigation/elements';
import {SafeAreaProvider, SafeAreaView} from "react-native-safe-area-context";

import React, {useEffect, useState} from 'react';
import {Text, TouchableOpacity, View} from 'react-native';
import {Link} from "expo-router";

function showAlert(message: string) {
  if (Platform.OS === 'web') {
    window.alert(message);
  } else {
    Alert.alert(message);
  }
}
export default function TabTwoScreen() {
    const [isOn, setIsOn] = useState(false);
    const onPress = () => setIsOn(prev => !prev);


    const [devices, setDevices] = useState<any[]>([]);
    useEffect(() => {
        getDevices();
    }, []);
    const getDevices = async () => {
        try {
            const login = await fetch("http://localhost:3000/auth/login", {
                method: "POST",
                headers: {"Content-Type": "application/json",},
                body: JSON.stringify({
                    accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzM1NDQ5MzMsImlhdCI6MTc3MzUzNzczMywiaXNzIjoiZGVjaXNpb24tbWFrZXIiLCJzdWIiOiJ0QjJXbUZ4eXIzcHQifQ.halo4p0QLYRXsroV5MANlq5JyZAb7W792sptdPeqGlI",
                    refreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzYyMTYxMzMsImlhdCI6MTc3MzUzNzczMywiaXNzIjoiZGVjaXNpb24tbWFrZXIiLCJzdWIiOiJ0QjJXbUZ4eXIzcHQifQ.HM2IFo5ks7R3SlFMf556XiNe-BK0tV2-rukaz2cSD1g",
                    userId: "tB2WmFxyr3pt",
                    username: "caiocotts"
                }),
            });
            const loginData = await login.json();
            const token = loginData.accessToken;
            const response = await fetch("http://localhost:3000/api/device", {
                method: "GET",
                headers: {"Content-Type": "application/json", Authorization: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzM1NDQ5MzMsImlhdCI6MTc3MzUzNzczMywiaXNzIjoiZGVjaXNpb24tbWFrZXIiLCJzdWIiOiJ0QjJXbUZ4eXIzcHQifQ.halo4p0QLYRXsroV5MANlq5JyZAb7W792sptdPeqGlI',},
            });
            const data = await response.json();
            setDevices(data);
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

        <View>

            <TouchableOpacity style={styles.button} onPress={getDevices}>
                <Text style={styles.addButtonText}>Get Devices</Text>
            </TouchableOpacity>

            {devices.map((device, index) => (
                <Text key={index}>{JSON.stringify(device)}</Text>
            ))}
        </View>
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

        {/* end of new */}


        <View style={styles.row}>
                <View style={styles.countContainer}>
                    <Text style={{ color: '#FFFFFF' }}>Device 1: {isOn ? 'ON' : 'OFF'}</Text>
                </View>
                <TouchableOpacity style={styles.button} onPress={onPress}>
                    <Text>Turn {isOn ? 'OFF' : 'ON'} </Text>
                </TouchableOpacity>
        </View>


        <View style={styles.row}>
            <View style={styles.countContainer}>
                <Text style={{ color: '#FFFFFF' }}>Device 1: {isOn ? 'ON' : 'OFF'}</Text>
            </View>
            <TouchableOpacity style={styles.button} onPress={onPress}>
                <Text>Turn {isOn ? 'OFF' : 'ON'} </Text>
            </TouchableOpacity>
        </View>
        {/* end of new */}


        <Link href="/modal" asChild>
            <TouchableOpacity style={styles.addButton}>
                <Text style={styles.addButtonText}>Add device</Text>
            </TouchableOpacity>
        </Link>


    </ParallaxScrollView>
  );
}

const styles = StyleSheet.create({
    addButton: {
        backgroundColor: '#4e8cc2',
        padding: 10,
        borderRadius: 8,
        alignItems: 'center',
        width: 200
    },

    addButtonText: {
        fontSize: 16,
    },
    row: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: 10, // space between items
    },
    container: {
        flex: 1,
        justifyContent: 'center',
        paddingHorizontal: 10,
    },
    button: {
        alignItems: 'center',
        backgroundColor: '#DDDDDD',
        padding: 10,
        width: 100,
    },
    countContainer: {
        alignItems: 'center',
        padding: 10,
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
