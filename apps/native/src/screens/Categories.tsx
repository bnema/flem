import React from 'react'
import { StyleSheet, View, Text } from 'react-native'
import { COLORS } from '../utils/Params'
import Genres from '../components/genres/Genres'

export default function Categories() {
    return (
        <View style={styles.container} >
            <Genres />
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: COLORS.BACKGROUND,
    }
})