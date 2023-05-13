import { StatusBar } from 'expo-status-bar'
import { NavigationContainer } from '@react-navigation/native'
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs'
import FontAwesome from 'react-native-vector-icons/FontAwesome'
import { COLORS } from './src/utils/Params'
import Discover from './src/screens/Discover'
import MySelection from './src/screens/MySelection'

const Tab = createBottomTabNavigator()

export default function App() {
  return (
    <>
      <NavigationContainer>
        <Tab.Navigator screenOptions={({ route }) => ({
          tabBarIcon: ({ color, size }) => {
            let iconName = ''

            if (route.name === 'Discover') {
              iconName = 'search'
            } else if (route.name === 'MySelection') {
              iconName = 'heart'
            }

            return <FontAwesome name={iconName} size={size} color={color} />;
          },
          tabBarActiveTintColor: COLORS.ICON,
          tabBarInactiveTintColor: '#737373',
          tabBarStyle: {
            height: 72,
            paddingHorizontal: 7,
            paddingBottom: 7,
            borderTopWidth: 0,
            backgroundColor: COLORS.UI,
          },
          headerShown: false,
        })}>
          <Tab.Screen name='Discover' component={Discover} />
          <Tab.Screen name='MySelection' component={MySelection} />
        </Tab.Navigator>
      </NavigationContainer>
      <StatusBar style='auto' />
    </>
  )
}