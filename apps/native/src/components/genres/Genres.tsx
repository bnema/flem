import React, { useEffect, useState } from 'react'
import { StyleSheet, View, Text } from 'react-native'
import GenreSelection from '../selection/GenreSelection'
import genresDB from '../../utils/genresDB'
import { COLORS } from '../../utils/Params'

export default function Genres() {


    return (
        <View>
            <Text style={styles.text}>What should you watch?</Text>
            <Text style={styles.note}>Choose as many as you want</Text>
            <View style={styles.genresContainer}>
                {genresDB.map(({ id, name }) => {
                    return <GenreSelection key={id} id={id} genre={name} />
                })}
            </View>
        </View>
    )
}

const styles = StyleSheet.create({
    text: {
        paddingTop: 30,
        left: 20,
        fontSize: 20,
        fontWeight: '500',
        color: COLORS.VIEW,
        marginBottom: 5,
    },
    note: {
        left: 20,
        fontSize: 14,
        fontWeight: '300',
        color: COLORS.VIEW,
        marginBottom: 30,
    },
    genresContainer: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        columnGap: 15,
        rowGap: 20,
        paddingHorizontal: 20,
    }
})